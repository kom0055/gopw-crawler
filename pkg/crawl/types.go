package crawl

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

type Encodee interface {
	GetSrvCode() SrvCode
	GetChannelCode() ChannelCode
	GetReqType() ReqType
}

func NewRequest[T Encodee](data T) *Request[T] {
	return &Request[T]{
		Data: data,
	}
}

type Request[T Encodee] struct {
	SrvCode     SrvCode     `json:"srvCode"`
	SignData    string      `json:"signData"`
	Sign        string      `json:"sign"`
	ChannelCode ChannelCode `json:"channelCode"`
	Timestamp   int64       `json:"timestamp"`
	Data        T           `json:"-"`
}

func (r *Request[T]) MarshalJSON() ([]byte, error) {
	sign, signData, err := TransEncodee(r.Data)
	if err != nil {
		return nil, err
	}
	r.SignData = signData
	r.Sign = sign
	r.ChannelCode = r.Data.GetChannelCode()
	r.Timestamp = time.Now().Unix()
	r.SrvCode = r.Data.GetSrvCode()
	type plain Request[T]
	pp := plain(*r)
	b, err := json.Marshal(pp)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func TransEncodee(data Encodee) (string, string, error) {
	paramBytes, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}
	paramStr := url.QueryEscape(string(paramBytes))
	signDataBytes := []byte(paramStr)
	switch data.GetReqType() {
	case ReqTypeEmpty, ReqType01, ReqType02, ReqTypeIsG1:
		pubKeyHexCopy := make([]byte, len(pubkeyHex))
		copy(pubKeyHexCopy, pubkeyHex)
		//if len(pubKeyHexCopy) > 128 {
		//	pubKeyHexCopy = pubKeyHexCopy[:128]
		//}
		pubKey, err := x509.ReadPublicKeyFromHex(string(pubKeyHexCopy))
		if err != nil {
			return "", "", err
		}
		signDataBytes, err = sm2.Encrypt(pubKey, []byte(base64.StdEncoding.EncodeToString(signDataBytes)), rand.Reader, sm2.C1C2C3)
		if err != nil {
			return "", "", err
		}
	}
	hasher := md5.New()
	hasher.Write(signDataBytes)
	sign := hasher.Sum(nil)
	return hex.EncodeToString(sign), hex.EncodeToString(signDataBytes), nil
}

type BizParam struct {
	SrvCode     SrvCode     `json:"srvCode"`
	LogId       string      `json:"logId"`
	SerialNo    string      `json:"serialNo"`
	ChannelCode ChannelCode `json:"channelCode"`
	FuncCode    FuncCode    `json:"funcCode"`
	IsToken     IsToken     `json:"isToken"`
	IsBindSell  IsBindSell  `json:"isBindSell"`
	HasRight    HasRight    `json:"hasRight"`
}

func NewLoginParam(mobile, randomCode, passwd string) *LoginParam {
	return &LoginParam{
		BizParam: BizParam{
			SrvCode:     LoginSrvCode,
			LogId:       uuid.NewString(),
			SerialNo:    time.Now().Format(timeFmt),
			ChannelCode: ChannelCodeGOPW,
			FuncCode:    FuncCode54,
			IsToken:     IsToken0,
			IsBindSell:  IsBindSellP1,
			HasRight:    HasRightFalse,
		},
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

func NewLogoutParam(mobile string) *LogoutParam {
	return &LogoutParam{
		BizParam: BizParam{
			SrvCode:     LogoutSrvCode,
			LogId:       uuid.NewString(),
			SerialNo:    time.Now().Format(timeFmt),
			ChannelCode: ChannelCodeGOPW,
			FuncCode:    FuncCode54,
			IsToken:     IsToken0,
			IsBindSell:  IsBindSellP1,
			HasRight:    HasRightFalse,
		},
		MobileNo:    mobile,
		LoginMobile: mobile,
	}
}

type LogoutParam struct {
	BizParam    `json:",inline"`
	MobileNo    string `json:"mobileNo"`
	LoginMobile string `json:"loginMobile"`
}

type Response[T any] struct {
	ResponseText string `json:"responseText"`
	Data         T      `json:"-"`
}

func (r *Response[T]) UnmarshalJSON(crypted []byte) error {
	keyByte := []byte(desDecryptKey)
	type plain Response[T]
	pp := plain(*r)
	if err := json.Unmarshal(crypted, &pp); err != nil {
		return err
	}
	decodedStr, err := base64.StdEncoding.DecodeString(pp.ResponseText)
	if err != nil {
		return err
	}
	out, err := EcbDesDecrypt(decodedStr, keyByte)
	if err != nil {
		return err
	}
	var data T
	if err := json.Unmarshal(out, &data); err != nil {
		return err
	}
	r.Data = data
	return nil

}

type BizResp struct {
	RtnMsg  string  `json:"rtnMsg"`
	RtnCode RtnCode `json:"rtnCode"`
}

type LogoutResp struct {
	BizResp `json:",inline"`
}
