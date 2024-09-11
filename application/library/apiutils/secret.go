package apiutils

import (
	"net/url"
	"time"

	"github.com/coscms/webcore/library/common"
	"github.com/webx-top/com"
)

func Decrypt(decrypted string, secret string) string {
	raw := decrypted
	if len(raw) > 0 {
		raw = com.URLSafeBase64(raw, false)
	}
	common.Decrypt(secret, &raw)
	return raw
}

func Encrypt(raw string, secret string) string {
	encrypted := raw
	common.Encrypt(secret, &encrypted)
	if len(encrypted) > 0 {
		encrypted = com.URLSafeBase64(encrypted, true)
	}
	return encrypted
}

// BuildURLValues 构建API参数值
func BuildURLValues(values url.Values, secret string) url.Values {
	if values == nil {
		values = url.Values{}
	}
	//values.Set(`appID`, appID)
	values.Set(`timestamp`, com.ToStr(time.Now().Unix()))
	values.Del(`sign`)
	sign := MakeSign(values, secret)
	values.Set(`sign`, sign)
	return values
}

// MakeSign 生成签名
// 规则说明：数据按照以键的字母顺序排列后进行url编码，然后拼接上`&secret=密钥`
// 等价的PHP代码：
// ```php
// <?php
//
//	function makeSign(array $data, string $secret): string{
//		ksort($data);
//		return hash('sha256',http_build_query($data).'&secret='.$secret);
//	}
//
// ```
func MakeSign(data url.Values, secret string) string {
	return com.Sha256(data.Encode() + `&secret=` + secret)
}
