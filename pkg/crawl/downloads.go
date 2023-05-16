package crawl

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

func NewDownloadParam(loginInfo *LoginResp, queryDate time.Time, consRq ConsPq, logId string) Encodee {

	y, m, _ := queryDate.Date()
	startDate := time.Date(y, m, 0, 0, 0, 0, 0, time.Local)
	return &DownloadParam{
		BizParam:     NewBizParam(QuerySrvCode, ChannelCodeGOPW, FuncCode54, IsToken1, IsBindSell1, HasRightTrue, logId),
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
			BeginDate:   startDate.Format("2006-01-02"),
			MgtOrgCode:  consRq.OrgNo,
			CustNo:      consRq.ConsNo,
			EndDate:     queryDate.Format("2006-01-02"),
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
	EspFlowId    string `yaml:"espFlowId"`
	Token        string `yaml:"token"`
	MobileNo     string `yaml:"mobileNo"`
	LoginMobile  string `yaml:"loginMobile"`
	EspRsvField1 string `yaml:"espRsvField1"`
	EspRsvField2 string `yaml:"espRsvField2"`
	EspRsvField3 string `yaml:"espRsvField3"`
	EspTimestamp int64  `yaml:"espTimestamp"`
	EspSign      string `yaml:"espSign"`

	PageNo   string `json:"pageNo"`
	PageSize string `json:"pageSize"`

	EspInformation EspInformation `yaml:"espInformation"`
	Header         Header         `yaml:"header"`
}

type EspInformation struct {
	DevType     string `yaml:"devType"`
	BeginDate   string `yaml:"beginDate"`
	MgtOrgCode  string `yaml:"mgtOrgCode"`
	CustNo      string `yaml:"custNo"`
	EndDate     string `yaml:"endDate"`
	MeasAssetNo string `yaml:"measAssetNo"`
	DisplayType string `yaml:"displayType"`
	QueryType   string `yaml:"queryType"`
	TmnlAssetNo string `yaml:"tmnlAssetNo"`
}

type Header struct {
	Bizcode    string `yaml:"bizcode"`
	Reptag     string `yaml:"reptag"`
	Accesscode string `yaml:"accesscode"`
	Sessionid  string `yaml:"sessionid"`
}

func (l DownloadParam) GetReqType() ReqType {
	return ReqType01
}

type DownloadOption struct {
	LoginInfo *LoginResp
	QueryDate time.Time
	ConsRq    ConsPq
	LogId     string
}
