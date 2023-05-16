package crawl

func NewGetCodeParam(logId string) Encodee {
	return GetCodeParam{
		BizParam: NewBizParam(GetCodeSrvCode, ChannelCodeGOPW, FuncCode54, IsToken0, IsBindSellN1, HasRightFalse, logId),
	}

}

type GetCodeParam struct {
	BizParam `json:",inline"`
}

func (g GetCodeParam) GetReqType() ReqType {
	return ReqType01
}

type GetCodeResp struct {
	BizResp    `json:",inline"`
	RandomCode string `json:"randomCode"`
}
