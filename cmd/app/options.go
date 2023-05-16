package app

import (
	"context"
	"fmt"
	"gopw-crawler/pkg/crawl"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

const (
	dateFmt = "2006-01"
)

func NewOptions() *Options {
	return &Options{
		Mobile:       "",
		PassWD:       "",
		DownloadPath: "",
	}
}

type Options struct {
	Mobile       string    `yaml:"mobile"`
	PassWD       string    `yaml:"passWD"`
	DownloadPath string    `yaml:"downloadPath"`
	StartDate    string    `yaml:"startDate"`
	EndDate      string    `yaml:"endDate"`
	Start        time.Time `yaml:"-"`
	End          time.Time `yaml:"-"`
}

func (o *Options) Init() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if b, err := os.ReadFile(fmt.Sprintf("%s/%s", pwd, "crawler.yaml")); err == nil {
		if err := yaml.Unmarshal(b, o); err != nil {
			return err
		}
	}

	if len(o.DownloadPath) == 0 {
		o.DownloadPath = pwd
	}
	return nil
}

func (o *Options) AddAllFlags(fs *pflag.FlagSet) {

	fs.StringVar(&o.Mobile, "mobile", o.Mobile, "mobile")
	fs.StringVar(&o.PassWD, "passwd", o.PassWD, "password")
	fs.StringVar(&o.DownloadPath, "download-path", o.DownloadPath, "download path")
	fs.StringVar(&o.StartDate, "start-date", o.StartDate, "start date")
	fs.StringVar(&o.EndDate, "end-date", o.EndDate, "end date")
}

func (o *Options) Complete() error {
	if len(o.Mobile) == 0 {
		return fmt.Errorf("mobile is empty")
	}
	if len(o.PassWD) == 0 {
		return fmt.Errorf("password is empty")
	}
	encryptedPasswd, err := crawl.EncryptPasswd(o.PassWD)
	if err != nil {
		return err
	}
	o.PassWD = encryptedPasswd
	now := time.Now()
	y, m, _ := now.Date()
	o.Start = time.Date(y, m, 0, 0, 0, 0, 0, time.Local)
	if len(o.StartDate) > 0 {
		start, err := time.Parse(dateFmt, o.StartDate)
		if err == nil {
			o.Start = start
		}
	}
	o.End = o.Start.AddDate(0, -6, 0)
	if len(o.EndDate) > 0 {
		end, err := time.Parse(dateFmt, o.EndDate)
		if err == nil {
			o.End = end
		}
	}

	return nil
}

func (o *Options) Run() error {
	return ServeWithSignalContext(o.run)
}
func (o *Options) run(ctx context.Context) error {
	downloadPath := fmt.Sprintf("%s/%s", o.DownloadPath, "gopw-download")
	if err := os.MkdirAll(downloadPath, 0777); err != nil {
		return err
	}
	crawler := crawl.DailyPowerCrawler{
		DownloadPath: downloadPath,
		Mobile:       o.Mobile,
		Passwd:       o.PassWD,
	}
	if err := crawler.Crawl(ctx, o.Start, o.End); err != nil {
		return err
	}
	return nil
}

// ServeWithSignalContext is a helper function that runs a server until a signal is received.
func ServeWithSignalContext(srv func(ctx context.Context) error) error {
	// shutdownSignals is the list of signals that can trigger a graceful shutdown.
	shutdownSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, os.Kill}
	shutdownChannel := make(chan os.Signal, 2)

	signal.Notify(shutdownChannel, shutdownSignals...)
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error)
	defer close(errCh)
	go func() {
		defer cancel()
		if err := srv(ctx); err != nil {
			errCh <- err
		}
	}()
	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	case sig := <-shutdownChannel:
		log.Printf("receive first signal. sig: %v", sig)
		cancel()
		sig = <-shutdownChannel
		log.Printf("receive second signal. sig: %v", sig)
	}
	return nil

}
