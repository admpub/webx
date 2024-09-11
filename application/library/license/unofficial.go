//go:build !official
// +build !official

package license

import (
	"github.com/coscms/webcore/library/license"
)

func init() {
	license.SkipLicenseCheck = false
	/*
		(&license.ServerURL{
			Tracker: `http://127.0.0.1:8080/product/script/nging/tracker.js`,
			Product: `http://127.0.0.1:8080/product/detail/nging`,
			License: `http://127.0.0.1:8080/product/license/nging`,
		}).Apply()
	// */
}
