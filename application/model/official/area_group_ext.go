package official

import "github.com/admpub/webx/application/dbschema"

type AreaGroupExt struct {
	*dbschema.OfficialCommonAreaGroup
	Areas []*dbschema.OfficialCommonArea `db:"-,relation=id:area_ids|notEmpty|split"`
}
