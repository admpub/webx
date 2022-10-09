package top

import (
	"regexp"
	"strings"
)

var (
	mobileHeader = regexp.MustCompile(`\b(?i:Mobile|iPod|iPad|iPhone|Android|Opera Mini|BlackBerry|webOS|UCWEB|Blazer|PSP)\b`)
)

func IsMobile(userAgent string) bool {
	return mobileHeader.MatchString(userAgent)
}

func IsWechat(userAgent string) bool {
	return strings.Contains(userAgent, `MicroMessenger`)
}
