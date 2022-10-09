package official

import "github.com/admpub/webx/application/dbschema"

type ICategory interface {
	GetCategory1() uint
	GetCategory2() uint
	GetCategory3() uint
	GetCategoryID() uint
	AddCategory(*dbschema.OfficialCommonCategory)
}
