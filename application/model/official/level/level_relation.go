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
