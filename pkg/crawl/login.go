package crawl

func NewLoginParam(mobile, passwd, randomCode, logId string) Encodee {
	return &LoginParam{
		BizParam:    NewBizParam(LoginSrvCode, ChannelCodeGOPW, FuncCode54, IsToken0, IsBindSellN1, HasRightFalse, logId),
		MobileNo:    mobile,
		LoginMobile: mobile,
		RandomCode:  randomCode,
		PassWord:    passwd,
	}
}

type LoginParam struct {
	BizParam    `json:",inline"`
	MobileNo    string `json:"mobileNo"`
	LoginMobile string `json:"loginMobile"`
	RandomCode  string `json:"randomCode"`
	PassWord    string `json:"passWord"`
}

func (l LoginParam) GetReqType() ReqType {
	return ReqType01
}

type LoginResp struct {
	BizResp      `json:",inline"`
	CertType     string     `json:"certType"`
	AuthenStatus string     `json:"authenStatus"`
	IsBindGroup  string     `json:"isBindGroup"`
	MobileNo     string     `json:"mobileNo"`
	IsBindSell   string     `json:"isBindSell"`
	Token        string     `json:"token"`
	CertNo       string     `json:"certNo"`
	AccountId    string     `json:"accountId"`
	LogFailNum   int        `json:"logFailNum"`
	CertName     string     `json:"certName"`
	IsBindCons   string     `json:"isBindCons"`
	ClientIp     string     `json:"clientIp"`
	SellList     []SellInfo `json:"sellList"`
}

type SellInfo struct {
	CorpName       string `json:"corpName"`
	CorpNo         string `json:"corpNo"`
	UnionCreditNo  string `json:"unionCreditNo"`
	BindId         string `json:"bindId"`
	BindTelNo      string `json:"bindTelNo"`
	LegalName      string `json:"legalName"`
	AccountId      string `json:"accountId"`
	CorpRegAddr    string `json:"corpRegAddr"`
	IsDefault      string `json:"isDefault"`
	BindTime       int64  `json:"bindTime"`
	LegalContactNo string `json:"legalContactNo"`
	State          string `json:"state"`
	BindWay        string `json:"bindWay"`
	BindTypeCode   string `json:"bindTypeCode"`
	ChannelCode    string `json:"channelCode"`
}
