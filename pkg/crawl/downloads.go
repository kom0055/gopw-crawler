package crawl

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

func NewDownloadParam(loginInfo *LoginResp, queryDate time.Time, consRq ConsPq, logId string) Encodee {

	y, m, _ := queryDate.Date()
	end := time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
	end = end.AddDate(0, 1, -1)
	return &DownloadParam{
		BizParam:     NewBizParam(DownloadSrcCode, ChannelCodeGOPW, FuncCode54, IsToken1, IsBindSell1, HasRightTrue, logId),
		EspFlowId:    strings.ReplaceAll(uuid.NewString(), "-", ""),
		Token:        loginInfo.Token,
		MobileNo:     loginInfo.MobileNo,
		LoginMobile:  loginInfo.MobileNo,
		EspRsvField1: "",
		EspRsvField2: "",
		EspRsvField3: "",
		EspTimestamp: time.Now().Unix(),
		EspSign:      "",
		PageNo:       "",
		PageSize:     "20000",
		EspInformation: EspInformation{
			DevType:     "02",
			BeginDate:   queryDate.Format("2006-01-02"),
			MgtOrgCode:  consRq.OrgNo,
			CustNo:      consRq.ConsNo,
			EndDate:     end.Format("2006-01-02"),
			MeasAssetNo: "",
			DisplayType: "",
			QueryType:   "",
			TmnlAssetNo: "",
		},
		Header: Header{
			Bizcode:    "2010010057",
			Reptag:     "0",
			Accesscode: "201001",
			Sessionid:  "201001005720220420000000100800",
		},
	}
}

type DownloadParam struct {
	BizParam     `json:",inline"`
	EspFlowId    string `json:"espFlowId"`
	Token        string `json:"token"`
	MobileNo     string `json:"mobileNo"`
	LoginMobile  string `json:"loginMobile"`
	EspRsvField1 string `json:"espRsvField1"`
	EspRsvField2 string `json:"espRsvField2"`
	EspRsvField3 string `json:"espRsvField3"`
	EspTimestamp int64  `json:"espTimestamp"`
	EspSign      string `json:"espSign"`

	PageNo   string `json:"pageNo"`
	PageSize string `json:"pageSize"`

	EspInformation EspInformation `json:"espInformation"`
	Header         Header         `json:"header"`
}

type EspInformation struct {
	DevType     string `json:"devType"`
	BeginDate   string `json:"beginDate"`
	MgtOrgCode  string `json:"mgtOrgCode"`
	CustNo      string `json:"custNo"`
	EndDate     string `json:"endDate"`
	MeasAssetNo string `json:"measAssetNo"`
	DisplayType string `json:"displayType"`
	QueryType   string `json:"queryType"`
	TmnlAssetNo string `json:"tmnlAssetNo"`
}

type Header struct {
	Bizcode    string `json:"bizcode"`
	Reptag     string `json:"reptag"`
	Accesscode string `json:"accesscode"`
	Sessionid  string `json:"sessionid"`
}

func (l DownloadParam) GetReqType() ReqType {
	return ReqTypeNoEncrypt
}

type DownloadOption struct {
	FileName  string
	LoginInfo *LoginResp
	QueryDate time.Time
	ConsRq    ConsPq
	LogId     string
}
