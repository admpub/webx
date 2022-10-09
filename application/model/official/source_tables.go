package official

import (
	"github.com/webx-top/echo"
)

// SourceTableEntry 来源表项
type SourceTableEntry struct {
	// Name 名称
	Name string

	// QueryKVAndDetailURL 查询id和名称以及详情页网址
	// H: echo.H{
	//	`detailURL`: `%v`,
	//},
	QueryKVAndDetailURL func(ctx echo.Context, sourceIDs []string) echo.KVList

	// QueryBought 查询是否购买过
	QueryBought func(ctx echo.Context, sourceID string, customerID []uint64) (map[uint64]bool, error)
}

// SourceTables 来源表集合
// used: comment
var SourceTables = map[string]*SourceTableEntry{}
