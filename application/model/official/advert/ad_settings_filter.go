package advert

import "github.com/webx-top/db"

type Filter struct {
	Client string
	AreaID uint
	Age    uint
	NowTS  int64
	Gendar string
}

func (f *Filter) GenCond() *db.Compounds {
	cond := db.NewCompounds()
	if len(f.Client) > 0 {
		cond.Add(db.And(
			db.Cond{`type`: `client`},
			db.Cond{`value`: db.In([]string{f.Client, ``})},
		))
	}
	if f.Age > 0 {
		cond.Add(db.And(
			db.Cond{`type`: `age`},
			db.Or(
				db.Cond{`v_start`: 0},
				db.Cond{`v_start`: db.Gte(f.Age)},
			),
			db.Or(
				db.Cond{`v_end`: 0},
				db.Cond{`v_end`: db.Lte(f.Age)},
			),
		))
	}
	if f.NowTS > 0 {
		cond.Add(db.And(
			db.Cond{`type`: `time`},
			db.Or(
				db.Cond{`v_start`: 0},
				db.Cond{`v_start`: db.Gte(f.NowTS)},
			),
			db.Or(
				db.Cond{`v_end`: 0},
				db.Cond{`v_end`: db.Lte(f.NowTS)},
			),
		))
	}
	if len(f.Gendar) > 0 {
		cond.Add(db.And(
			db.Cond{`type`: `gendar`},
			db.Cond{`value`: db.In([]string{f.Gendar, ``})},
		))
	}
	return cond
}
