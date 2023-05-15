package crawl

import (
	"time"

	"github.com/google/uuid"
)

func NewGetCodeParam() Encodee {
	return GetCodeParam{
		BizParam: BizParam{
			SrvCode:     GetCodeSrvCode,
			LogId:       uuid.NewString(),
			SerialNo:    time.Now().Format(timeFmt),
			ChannelCode: ChannelCodeGOPW,
			FuncCode:    FuncCode54,
			IsToken:     IsToken0,
			IsBindSell:  IsBindSellP1,
			HasRight:    HasRightFalse,
		},
	}
}

type GetCodeParam struct {
	BizParam `json:",inline"`
}

func (g GetCodeParam) GetReqType() ReqType {
	return ReqType01
}

func (g GetCodeParam) GetSrvCode() SrvCode {
	return GetCodeSrvCode
}

func (g GetCodeParam) GetChannelCode() ChannelCode {
	return ChannelCodeGOPW
}

type GetCodeResp struct {
	BizResp    `json:",inline"`
	RandomCode string `json:"randomCode"`
}
