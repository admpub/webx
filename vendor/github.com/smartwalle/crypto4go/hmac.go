package crypto4go

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func HmacMD5(data, key []byte) []byte {
	var h = hmac.New(md5.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func HmacMD5String(data, key string) string {
	return hex.EncodeToString(HmacMD5([]byte(data), []byte(key)))
}

func HmacSHA1(data, key []byte) []byte {
	var h = hmac.New(sha1.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func HmacSHA1String(data, key string) string {
	return hex.EncodeToString(HmacSHA1([]byte(data), []byte(key)))
}

func HmacSHA256(data, key []byte) []byte {
	var h = hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func HmacSHA256String(data, key string) string {
	return hex.EncodeToString(HmacSHA256([]byte(data), []byte(key)))
}

func HmacSHA512(data, key []byte) []byte {
	var h = hmac.New(sha512.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func HmacSHA512String(data, key string) string {
	return hex.EncodeToString(HmacSHA512([]byte(data), []byte(key)))
}
