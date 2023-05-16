package crawl

const (
	pubkeyHex = "0419ebd19bd923575a88d06bd6715e8de98c4d11f8a73fc5a98421ddae54ce750280c2b0cbdd6eaf224dac1dfc042bf4918b8171c61e50a526fbfadbd9068a8510"
	publicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCNaLDNXeTB29zaFvYorSTaE+Ux
XmLAXiG7qc8ExLo4+GCHqoNrW2GbWn1RxUj98BgdWoyFu91nTvDPs/3vF9xZw9zZ
pjE+budWH+9Fm349DlivzDD2LAbmJcTeR96X6VjAaGyLYdIWlue4cgG9ZlDKGDvc
/8r5SKzveRtTt18XSQIDAQAB
-----END PUBLIC KEY-----`
	timeFmt       = "20060102150405"
	desDecryptKey = "60d1baf4d4034239bcfec7d321dee794"
)

type SrvCode string

const (
	GetCodeSrvCode  = "LSSP_002381"
	LoginSrvCode    = "LSSP_001048"
	LogoutSrvCode   = "LSSP_001382"
	QuerySrvCode    = "LSSP_003108"
	DownloadSrcCode = "EMSS_001506"
)

type ChannelCode string

const (
	ChannelCodeGOPW ChannelCode = "GOPW"
)

type ReqType string

const (
	ReqTypeEmpty ReqType = ""
	ReqType01    ReqType = "01"
	ReqType02    ReqType = "02"
	ReqTypeIsG1          = "isg1"
)

type FuncCode string

const (
	FuncCode54 FuncCode = "54"
)

type IsToken string

const (
	IsToken0 IsToken = "0"
	IsToken1 IsToken = "1"
)

type IsBindSell string

const (
	IsBindSellN1 IsBindSell = "-1"
	IsBindSell1  IsBindSell = "1"
)

type HasRight string

const (
	HasRightFalse HasRight = "false"
	HasRightTrue  HasRight = "true"
)

type RtnCode string

const (
	RtnCodeOK RtnCode = "0"
)
