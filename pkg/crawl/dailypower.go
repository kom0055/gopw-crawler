package crawl

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"sync"
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
	client http.Client
	once   sync.Once
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

func (c *DailyPowerCrawler) Crawl(ctx context.Context) error {
	if err := c.Init(ctx); err != nil {
		return err
	}
	return nil
}

func (c *DailyPowerCrawler) GetCode(ctx context.Context) (*GetCodeResp, error) {
	param := NewRequest(NewGetCodeParam())
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

func WrapRequest(req *http.Request) {
	for k := range header {
		v := header[k]
		for i := range v {
			req.Header.Add(k, v[i])
		}
	}
}
