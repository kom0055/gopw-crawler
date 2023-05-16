package crawl

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"time"

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
	signData := url.QueryEscape(string(paramBytes))
	signDataBytes := []byte(signData)
	switch data.GetReqType() {
	case ReqTypeEmpty, ReqType01, ReqType02, ReqTypeIsG1:
		pubKeyHexCopy := make([]byte, len(pubkeyHex))
		copy(pubKeyHexCopy, pubkeyHex)
		pubKey, err := x509.ReadPublicKeyFromHex(string(pubKeyHexCopy))
		if err != nil {
			return "", "", err
		}
		signDataBytes, err = sm2.Encrypt(pubKey, []byte(base64.StdEncoding.EncodeToString(signDataBytes)), rand.Reader, sm2.C1C2C3)
		if err != nil {
			return "", "", err
		}
		signData = hex.EncodeToString(signDataBytes)
	}
	hasher := md5.New()
	hasher.Write(signDataBytes)
	sign := hasher.Sum(nil)
	return hex.EncodeToString(sign), signData, nil
}

func NewBizParam(srvCode SrvCode, chanCode ChannelCode, funcCode FuncCode, isToken IsToken, isBindSell IsBindSell,
	hasRight HasRight, logId string) BizParam {
	return BizParam{
		SrvCode:     srvCode,
		LogId:       logId,
		SerialNo:    time.Now().Format(timeFmt),
		ChannelCode: chanCode,
		FuncCode:    funcCode,
		IsToken:     isToken,
		IsBindSell:  isBindSell,
		HasRight:    hasRight,
	}
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

func (g BizParam) GetSrvCode() SrvCode {
	return g.SrvCode
}

func (g BizParam) GetChannelCode() ChannelCode {
	return g.ChannelCode
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

func (r BizResp) IsOK() bool {
	return r.RtnCode == RtnCodeOK
}
