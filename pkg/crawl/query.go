package crawl

import (
	"strconv"
	"time"
)

const (
	defaultPageNumber = 100
	queryDateFmt      = "2006-01-02"
)

func NewQueryParam(loginInfo *LoginResp, queryDate time.Time, page int, logId string) Encodee {
	pageStr := strconv.Itoa(page)
	var sellInfo SellInfo
	if len(loginInfo.SellList) > 0 {
		sellInfo = loginInfo.SellList[0]
	}
	defaultPageNumberStr := strconv.Itoa(defaultPageNumber)
	return &QueryParam{
		BizParam:     NewBizParam(QuerySrvCode, ChannelCodeGOPW, FuncCode54, IsToken1, IsBindSell1, HasRightTrue, logId),
		RegistUserNo: loginInfo.AccountId,
		SecNo:        sellInfo.CorpNo,
		CertNo:       loginInfo.CertNo,
		QueryDate:    queryDate.Format(queryDateFmt),
		ConsNo:       "",
		Page:         pageStr,
		Number:       defaultPageNumberStr,
		PageNo:       pageStr,
		PageSize:     defaultPageNumberStr,
		Token:        loginInfo.Token,
		Mobile:       loginInfo.MobileNo,
		MobileNo:     loginInfo.MobileNo,
		LoginMobile:  loginInfo.MobileNo,
	}
}

type QueryParam struct {
	BizParam     `json:",inline"`
	RegistUserNo string `json:"registUserNo"`
	SecNo        string `json:"secNo"`
	CertNo       string `json:"certNo"`
	QueryDate    string `json:"queryDate"`
	ConsNo       string `json:"consNo"`
	Page         string `json:"page"`
	Number       string `json:"number"`
	PageNo       string `json:"pageNo"`
	PageSize     string `json:"pageSize"`
	Token        string `json:"token"`
	Mobile       string `json:"mobile"`
	MobileNo     string `json:"mobileNo"`
	LoginMobile  string `json:"loginMobile"`
}

func (l QueryParam) GetReqType() ReqType {
	return ReqType01
}

type QueryResp struct {
	BizResp     `json:",inline"`
	Total       int      `json:"total"`
	RelaTraceID string   `json:"relaTraceID"`
	TotalNum    int      `json:"totalNum"`
	Status      string   `json:"status"`
	ConsPqList  []ConsPq `json:"consPqList"`
}

type ConsPq struct {
	GfPq        interface{} `json:"gfPq"`
	VoltCode    string      `json:"voltCode"`
	AmtYmd      string      `json:"amtYmd"`
	OrgNo       string      `json:"orgNo"`
	Tpq         interface{} `json:"tpq"`
	AuthEndDate string      `json:"authEndDate"`
	ConsName    string      `json:"consName"`
	DgPq        interface{} `json:"dgPq"`
	JfPq        interface{} `json:"jfPq"`
	ConsNo      string      `json:"consNo"`
	Voltage     string      `json:"voltage"`
	IsShow      string      `json:"isShow"`
}
