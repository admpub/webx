package binding

import (
	"time"

	"github.com/webx-top/echo"

	"github.com/admpub/log"
	dbschemaNging "github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/nging/v4/application/library/email"
)

func init() {
	email.AddCallback(func(cfg *email.Config, err error) {
		if cfg.ID == 0 {
			return
		}
		logM := &dbschemaNging.NgingSendingLog{}
		var status, result string
		if err != nil {
			status = `failure`
			result = err.Error()
		} else {
			status = `success`
			result = `发送成功`
		}
		err = logM.UpdateFields(nil, echo.H{
			`sent_at`: time.Now().Unix(),
			`result`:  result,
			`status`:  status,
		}, `id`, cfg.ID)
		if err != nil {
			log.Error(err)
		}
	})
}
