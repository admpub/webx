package page

import (
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
)

type LayoutExt struct {
	*dbschema.OfficialPageLayout
	Page  *dbschema.OfficialPage      `db:"-,relation=id:page_id|gtZero"`
	Block *dbschema.OfficialPageBlock `db:"-,relation=id:block_id|gtZero"`
}

type LayoutWithBlock struct {
	*dbschema.OfficialPageLayout
	Block       *dbschema.OfficialPageBlock `db:"-,relation=id:block_id|gtZero"`
	Configs     echo.H
	ItemConfigs echo.H
}
