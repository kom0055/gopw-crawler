package crawl

import (
	"bytes"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

func EcbDesEncrypt(origData, key []byte) ([]byte, error) {
	tkey := make([]byte, 8, 8)
	copy(tkey, key)

	block, err := des.NewCipher(tkey)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	origData = PKCS5Padding(origData, bs)

	out, err := encrypt(origData, tkey)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func EcbDesDecrypt(crypted, key []byte) ([]byte, error) {
	tkey := make([]byte, 8, 8)
	copy(tkey, key)

	out, err := decrypt(crypted, tkey)
	if err != nil {
		return nil, err
	}
	out = PKCS5Unpadding(out)
	return out, nil
}

// Des encryption
func encrypt(origData, key []byte) ([]byte, error) {
	if len(origData) < 1 || len(key) < 1 {
		return nil, errors.New("wrong data or key")
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(origData)%bs != 0 {
		return nil, errors.New("wrong padding")
	}
	out := make([]byte, len(origData))
	dst := out
	for len(origData) > 0 {
		block.Encrypt(dst, origData[:bs])
		origData = origData[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

// Des Decrypt
func decrypt(crypted, key []byte) ([]byte, error) {
	if len(crypted) < 1 || len(key) < 1 {
		return nil, errors.New("wrong data or key")
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(crypted))
	dst := out
	bs := block.BlockSize()
	if len(crypted)%bs != 0 {
		return nil, errors.New("wrong crypted size")
	}

	for len(crypted) > 0 {
		block.Decrypt(dst, crypted[:bs])
		crypted = crypted[bs:]
		dst = dst[bs:]
	}

	return out, nil
}

func PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func EncryptPasswd(passwd string) (string, error) {
	p, _ := pem.Decode([]byte(publicKey))
	if p == nil {
		return "", fmt.Errorf("public key error")
	}
	pubKeyAny, err := x509.ParsePKIXPublicKey(p.Bytes)
	if err != nil {
		return "", err
	}
	pubKey, ok := pubKeyAny.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("not ok")
	}
	passwdBytes := []byte(passwd)
	encryptedPasswd, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, passwdBytes)
	if err != nil {
		return "", err
	}
	encryptedPasswdStr := base64.StdEncoding.EncodeToString(encryptedPasswd)
	return encryptedPasswdStr, nil

}
