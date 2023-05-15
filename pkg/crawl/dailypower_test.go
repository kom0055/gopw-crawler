package crawl

import (
	"context"
	"testing"
)

func TestGetCode(t *testing.T) {
	crawler := DailyPowerCrawler{}
	ctx := context.Background()
	if err := crawler.Init(ctx); err != nil {
		t.Fatal(err)
	}
	resp, err := crawler.GetCode(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestTransGetCode(t *testing.T) {

	getCodeParam := GetCodeParam{
		BizParam: BizParam{
			SrvCode:     GetCodeSrvCode,
			LogId:       "b636fe5b-efe3-4fc9-842c-3e3c8200d78d",
			SerialNo:    "20230515233653",
			ChannelCode: "GOPW",
			FuncCode:    "54",
			IsToken:     "0",
			IsBindSell:  "-1",
			HasRight:    "false",
		},
	}
	sign, signData, err := TransEncodee(getCodeParam)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sign)
	t.Log(signData)

}
