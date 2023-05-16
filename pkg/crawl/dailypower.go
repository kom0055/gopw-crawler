package crawl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/http/cookiejar"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
)

const (
	logoutInterval = 10 * time.Second
	callInterval   = 500 * time.Millisecond
)

var (
	header = http.Header{
		"Accept":                    []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"Accept-Language":           []string{"zh-CN,zh;q=0.9"},
		"Connection":                []string{"keep-alive"},
		"Sec-Fetch-Dest":            []string{"document"},
		"Sec-Fetch-Mode":            []string{"navigate"},
		"Sec-Fetch-Site":            []string{"none"},
		"Sec-Fetch-User":            []string{"?1"},
		"Upgrade-Insecure-Requests": []string{"1"},
		"User-Agent":                []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
		"sec-ch-ua":                 []string{`"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`},
		"sec-ch-ua-mobile":          []string{"?0"},
		"sec-ch-ua-platform":        []string{`"macOS"`},
	}
	mainPageUrl = "https://zsdl.zj.sgcc.com.cn/isg/app/gopwfront/pages/main/main.html"
	apiUrl      = "https://zsdl.zj.sgcc.com.cn/zj_eqa/open/gopwInvoke"
)

type DailyPowerCrawler struct {
	client       http.Client
	once         sync.Once
	DownloadPath string
	Mobile       string
	Passwd       string
}

func (c *DailyPowerCrawler) Init(ctx context.Context) error {
	var (
		err error
	)
	c.once.Do(func() {
		var (
			jar http.CookieJar
			req *http.Request
		)
		jar, err = cookiejar.New(nil)
		if err != nil {
			return
		}
		c.client = http.Client{Jar: jar}
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, mainPageUrl, nil)
		if err != nil {
			return
		}
		WrapRequest(req)
		_, err = c.client.Do(req)
		return
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *DailyPowerCrawler) Crawl(ctx context.Context, startDate, endDate time.Time) error {
	if err := c.Init(ctx); err != nil {
		return err
	}
	logId := uuid.NewString()
	_, err := c.Logout(ctx, c.Mobile, logId)
	if err != nil {
		return err
	}
	time.Sleep(logoutInterval)
	getCodeResp, err := c.GetCode(ctx, logId)
	if err != nil {
		return err
	}
	if !getCodeResp.IsOK() {
		return fmt.Errorf(getCodeResp.RtnMsg)
	}
	randomCode := getCodeResp.RandomCode
	time.Sleep(callInterval)
	loginResp, err := c.Login(ctx, c.Mobile, c.Passwd, randomCode, logId)
	if err != nil {
		return err
	}
	if !loginResp.IsOK() {
		return fmt.Errorf(loginResp.RtnMsg)
	}
	defer func() {
		_, _ = c.Logout(ctx, c.Mobile, logId)

	}()
	end := endDate
	start := startDate

	downloadOpts := map[string]DownloadOption{}
	for cursor := end; !cursor.Before(start); cursor = cursor.AddDate(0, -1, 0) {
		for page, total := 1, 2; page < total; page++ {
			time.Sleep(callInterval)
			queryResp, err := c.Query(ctx, loginResp, cursor, page, logId)
			if err != nil {
				return err
			}
			if !queryResp.IsOK() {
				return fmt.Errorf(queryResp.RtnMsg)
			}
			total = int(math.Floor(float64(queryResp.TotalNum) / float64(defaultPageNumber)))
			for i := range queryResp.ConsPqList {
				consPq := queryResp.ConsPqList[i]
				d := DownloadOption{
					FileName:  fmt.Sprintf("%s/%s-%s-%s.xls", c.DownloadPath, cursor.Format("2006-01"), consPq.VoltCode, consPq.ConsName),
					LoginInfo: loginResp,
					QueryDate: cursor,
					ConsRq:    consPq,
					LogId:     logId,
				}
				downloadOpts[d.FileName] = d

			}

		}
	}
	log.Println("total", len(downloadOpts))
	done := atomic.Int64{}
	ctrlCh := make(chan struct{}, 8)
	errCh := make(chan error, 8)
	go func() {
		defer close(ctrlCh)
		defer close(errCh)
		wg := sync.WaitGroup{}
		for i := range downloadOpts {
			downloadOpt := downloadOpts[i]
			wg.Add(1)
			go func() {
				defer wg.Done()
				ctrlCh <- struct{}{}
				defer func() {
					<-ctrlCh
				}()
				if err := c.Download(ctx, downloadOpt); err != nil {
					errCh <- err
				}
				done.Add(1)
				log.Println("done", done.Load())
			}()
		}
		wg.Wait()

	}()
	var multiErr *multierror.Error
	for err := range errCh {
		multiErr = multierror.Append(multiErr, err)
	}

	return multiErr.ErrorOrNil()
}

func (c *DailyPowerCrawler) GetCode(ctx context.Context, logId string) (*GetCodeResp, error) {
	param := NewRequest(NewGetCodeParam(logId))
	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(param); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiUrl, buffer)
	if err != nil {
		return nil, err
	}
	WrapRequest(req)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	getCodeResp := &Response[GetCodeResp]{}
	if err := decoder.Decode(getCodeResp); err != nil {
		return nil, err
	}
	return &getCodeResp.Data, nil

}

func (c *DailyPowerCrawler) Login(ctx context.Context, mobile, passwd, randomCode, logId string) (*LoginResp, error) {
	param := NewRequest(NewLoginParam(mobile, passwd, randomCode, logId))
	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(param); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiUrl, buffer)
	if err != nil {
		return nil, err
	}
	WrapRequest(req)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	logineResp := &Response[LoginResp]{}
	if err := decoder.Decode(logineResp); err != nil {
		return nil, err
	}
	return &logineResp.Data, nil
}

func (c *DailyPowerCrawler) Logout(ctx context.Context, mobile, logId string) (*LogoutResp, error) {
	param := NewRequest(NewLogoutParam(mobile, logId))
	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(param); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiUrl, buffer)
	if err != nil {
		return nil, err
	}
	WrapRequest(req)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	logoutResp := &Response[LogoutResp]{}
	if err := decoder.Decode(logoutResp); err != nil {
		return nil, err
	}
	return &logoutResp.Data, nil
}

func (c *DailyPowerCrawler) Query(ctx context.Context, loginInfo *LoginResp, queryDate time.Time,
	page int, logId string) (*QueryResp, error) {
	param := NewRequest(NewQueryParam(loginInfo, queryDate, page, logId))
	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(param); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiUrl, buffer)
	if err != nil {
		return nil, err
	}
	WrapRequest(req)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	queryResp := &Response[QueryResp]{}
	if err := decoder.Decode(queryResp); err != nil {
		return nil, err
	}
	return &queryResp.Data, nil
}

func (c *DailyPowerCrawler) Download(ctx context.Context, downloadOpt DownloadOption) error {

	queryDate := downloadOpt.QueryDate
	consRq := downloadOpt.ConsRq
	loginInfo := downloadOpt.LoginInfo
	logId := downloadOpt.LogId
	fileName := downloadOpt.FileName
	param := NewRequest(NewDownloadParam(loginInfo, queryDate, consRq, logId))
	b, err := json.Marshal(param)
	if err != nil {
		return err
	}
	//tmpUrl := url.URL{Path: string(b)}
	//s := tmpUrl.EscapedPath()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		return err
	}
	//req.URL.RawQuery = fmt.Sprintf("jsonData=%v", s)
	q := req.URL.Query()
	q.Add("jsonData", string(b))
	req.URL.RawQuery = q.Encode()

	WrapRequest(req)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func WrapRequest(req *http.Request) {
	for k := range header {
		v := header[k]
		for i := range v {
			req.Header.Add(k, v[i])
		}
	}
}
