package top

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/echo/param"
)

func CondFindInSet(key string, value interface{}, useFulltextIndex ...bool) db.Compound {
	v := param.AsString(value)
	return mysql.FindInSet(key, v, useFulltextIndex...)
}
