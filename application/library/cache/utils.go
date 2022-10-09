package cache

import (
	"github.com/webx-top/com"
)

func IsDbAccount(v string) bool {
	if len(v) == 0 {
		return false
	}
	for _, r := range v {
		if !com.IsNumeric(r) {
			return false
		}
	}
	return true
}
