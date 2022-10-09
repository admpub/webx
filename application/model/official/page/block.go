package page

import (
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
)

func NewBlock(ctx echo.Context) *Block {
	m := &Block{
		OfficialPageBlock: dbschema.NewOfficialPageBlock(ctx),
	}
	return m
}

type Block struct {
	*dbschema.OfficialPageBlock
}
