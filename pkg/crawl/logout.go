package crawl

func NewLogoutParam(mobile, logId string) Encodee {
	return &LogoutParam{
		BizParam:    NewBizParam(LogoutSrvCode, ChannelCodeGOPW, FuncCode54, IsToken0, IsBindSellN1, HasRightFalse, logId),
		MobileNo:    mobile,
		LoginMobile: mobile,
	}
}

type LogoutParam struct {
	BizParam    `json:",inline"`
	MobileNo    string `json:"mobileNo"`
	LoginMobile string `json:"loginMobile"`
}

func (l LogoutParam) GetReqType() ReqType {
	return ReqType01
}

type LogoutResp struct {
	BizResp `json:",inline"`
}
