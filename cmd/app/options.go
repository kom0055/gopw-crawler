package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"
)

func NewOptions() *Options {

	return &Options{
		Mobile:       "",
		PassWD:       "",
		DownloadPath: "",
	}
}

type Options struct {
	Mobile       string
	PassWD       string
	DownloadPath string
}

func (o *Options) AddAllFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Mobile, "mobile", o.Mobile, "mobile")
	fs.StringVar(&o.PassWD, "passwd", o.PassWD, "password")
	fs.StringVar(&o.DownloadPath, "download-path", o.DownloadPath, "download path")
}

func (o *Options) Complete() error {
	if len(o.Mobile) == 0 {
		return fmt.Errorf("mobile is empty")
	}
	if len(o.PassWD) == 0 {
		return fmt.Errorf("password is empty")
	}
	if len(o.DownloadPath) == 0 {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		downloadPath := fmt.Sprintf("%s/%s", pwd, "downloads")
		if err = os.MkdirAll(downloadPath, 0644); err != nil {
			return err
		}
		o.DownloadPath = downloadPath
	}
	return nil
}

func (o *Options) Run() error {
	ServeWithSignalContext(o.run)
	return nil
}
func (o *Options) run(ctx context.Context) {

}

// ServeWithSignalContext is a helper function that runs a server until a signal is received.
func ServeWithSignalContext(srv func(ctx context.Context)) {
	// shutdownSignals is the list of signals that can trigger a graceful shutdown.
	shutdownSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, os.Kill}
	shutdownChannel := make(chan os.Signal, 2)

	signal.Notify(shutdownChannel, shutdownSignals...)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		srv(ctx)
	}()
	sig := <-shutdownChannel
	log.Printf("receive first signal. sig: %v", sig)
	cancel()
	sig = <-shutdownChannel
	log.Printf("receive second signal. sig: %v", sig)
}
