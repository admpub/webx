package license

import (
	"strings"
	"time"

	"github.com/admpub/license_gen/lib"
	"github.com/admpub/log"
	"github.com/coscms/webcore/library/license"
	"github.com/webx-top/com"
)

func init() {
	license.ValidateResultInitor = func() license.Validator {
		return &ValidateResult{}
	}
}

type ValidateResult struct {
	Date    string `json:"date"`
	Package string `json:"package"`
	Version string `json:"version"`
	Result  string `json:"result"`
	Signed  string `json:"signed"` // [Date][Version][Result]
}

// - 装包 -

func (v *ValidateResult) SignableBytes() []byte {
	return com.Str2bytes(`[` + v.Date + `][` + v.Package + `][` + v.Version + `][` + v.Result + `]`)
}

func (v *ValidateResult) SignData(privateKey string) (err error) {
	v.Signed, err = lib.SignData(privateKey, v.SignableBytes())
	return
}

// - 拆包 -

func (v *ValidateResult) UnsignData(publicKey string) (err error) {
	return lib.UnsignData(publicKey, v.Signed, v.SignableBytes())
}

func (v *ValidateResult) Validate() error {
	if v.Result != `OK` {
		return lib.ErrInvalidLicense
	}
	if v.Package != license.Package() {
		log.Warnf(`product package is mismatched: %v != %v`, v.Package, license.Package())
		return lib.ErrInvalidLicense
	}
	verNo := strings.SplitN(license.Version(), `-`, 2)[0]
	if !lib.CheckVersion(verNo, v.Version) {
		log.Warnf(`version number is mismatched: %v != %v`, v.Version, verNo)
		return lib.ErrInvalidLicense
	}

	timeUTC := time.Now().UTC()
	hours := timeUTC.Hour()
	if hours < 23 && hours > 1 {
		myDate := timeUTC.Format(`20060102`)
		if v.Date != myDate {
			log.Warnf(`system time is mismatched: %v != %v`, v.Date, myDate)
			return lib.ErrInvalidLicense
		}
	}
	publicKey := license.GetOrLoadPublicKey()
	if len(publicKey) == 0 {
		b, err := license.ReadLicenseKeyFile()
		if err != nil {
			return err
		}
		_, publicKey = license.LicenseDecode(b)
		if len(publicKey) > 0 {
			license.SetPublicKey(publicKey)
		} else {
			log.Warn(`license public key required`)
		}
	}
	return v.UnsignData(publicKey)
}
