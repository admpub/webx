package level

import "github.com/admpub/webx/application/dbschema"

type LevelGroup struct {
	Group string
	Title string
	List  []*dbschema.OfficialCustomerLevel
}
