package level

import "github.com/admpub/webx/application/dbschema"

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
