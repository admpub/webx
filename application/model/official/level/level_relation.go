package level

import (
	"time"

	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func NewRelation(ctx echo.Context) *Relation {
	m := &Relation{
		OfficialCustomerLevelRelation: dbschema.NewOfficialCustomerLevelRelation(ctx),
	}
	return m
}

// 客户的扩展组等级关联关系
type Relation struct {
	*dbschema.OfficialCustomerLevelRelation
}

func (f *Relation) ListByCustomerID(customerID uint64) ([]*dbschema.OfficialCustomerLevelRelation, error) {
	_, err := f.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`status`: `success`},
		db.Or(
			db.Cond{`expired`: 0},
			db.Cond{`expired`: db.Lt(time.Now().Unix())},
		),
	))
	if err != nil {
		return nil, err
	}
	return f.Objects(), nil
}

func (f *Relation) ListByCustomerIDs(customerIDs []uint64) (map[uint64][]*RelationExt, error) {
	list := []*RelationExt{}
	var mw func(db.Result) db.Result
	cond := db.NewCompounds()
	cond.Add(
		db.Cond{`customer_id`: db.In(customerIDs)},
		db.Cond{`status`: `success`},
		db.Or(
			db.Cond{`expired`: 0},
			db.Cond{`expired`: db.Lt(time.Now().Unix())},
		),
	)
	_, err := f.ListByOffset(&list, mw, 0, -1, db.And())
	if err != nil {
		return nil, err
	}
	results := map[uint64][]*RelationExt{}
	for _, row := range list {
		if _, ok := results[row.CustomerId]; !ok {
			results[row.CustomerId] = []*RelationExt{}
		}
		results[row.CustomerId] = append(results[row.CustomerId], row)
	}
	return results, nil
}
