package apiutils

import (
	"net/url"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/com"
)

// SignString 生成签名字符串
func SignString(raw string) string {
	return com.Md5(com.Md5(raw) + config.FromFile().APIKey())
}

// CheckSign 检查签名是否匹配
func CheckSign(raw string, sign string) error {
	if SignString(raw) != sign {
		return common.ErrInvalidSign
	}
	return nil
}

// GenSign 根据url.Values类型值生成签名
func GenSign(formData url.Values) string {
	formData.Del(`sign`)
	return SignString(formData.Encode())
}
