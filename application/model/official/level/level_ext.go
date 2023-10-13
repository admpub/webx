package level

import (
	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/echo"
)

const (
	LevelStatusActived = `actived`
	LevelStatusExpired = `expired`
)

var LevelStatuses = echo.NewKVData().Add(LevelStatusActived, `有效`).Add(LevelStatusExpired, `过期`)

type LevelGroup struct {
	Group string
	Title string
	List  []*dbschema.OfficialCustomerLevel
}

type RelationExt struct {
	*dbschema.OfficialCustomerLevelRelation
	Level *dbschema.OfficialCustomerLevel `db:"-,relation=id:level_id|gtZero"`
}

func (r *RelationExt) Name_() string {
	if r.OfficialCustomerLevelRelation == nil {
		r.OfficialCustomerLevelRelation = &dbschema.OfficialCustomerLevelRelation{}
	}
	return r.OfficialCustomerLevelRelation.Name_()
}
