package crypto4go

import (
	"bytes"
)

func ZeroPadding(data []byte, blockSize int) []byte {
	var diff = blockSize - len(data)%blockSize
	var paddingText = bytes.Repeat([]byte{0}, diff)
	return append(data, paddingText...)
}

func PKCS7Padding(data []byte, blockSize int) []byte {
	var diff = blockSize - len(data)%blockSize
	var paddingText = bytes.Repeat([]byte{byte(diff)}, diff)
	return append(data, paddingText...)
}

func PKCS7UnPadding(data []byte) []byte {
	var length = len(data)
	var unpadding = int(data[length-1])
	return data[:(length - unpadding)]
}
