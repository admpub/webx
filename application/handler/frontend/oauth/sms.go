package oauth

import (
	"github.com/coscms/sms"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/log"
	dbschemaNging "github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/registry/settings"
)

// SetSMSConfigs 第三方登录平台账号
func SetSMSConfigs() error {
	group := `sms`
	m := &dbschemaNging.NgingConfig{}
	_, err := m.ListByOffset(nil, nil, 0, 1, db.And(
		db.Cond{`group`: group},
		db.Cond{`disabled`: `N`},
	))
	//panic(com.Dump(m.Objects()))
	cfg := echo.H{}
	decoder := settings.GetDecoder(group)
	for _, row := range m.Objects() {
		if len(row.Value) == 0 {
			continue
		}
		cfg, err = settings.DecodeConfig(row, cfg, decoder)
		if err != nil {
			return err
		}
		sender, ok := cfg.GetStore(row.Key).Get(`ValueObject`).(sms.Sender)
		if !ok {
			continue
		}
		log.Debugf(`The SMS Provider is registered: %s`, row.Key)
		sms.Register(row.Key, sender)
	}
	return err
}

func UpdateSMSConfigs() error {
	sms.Clear()
	return SetSMSConfigs()
}
