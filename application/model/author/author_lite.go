package author

import "github.com/admpub/webx/application/dbschema"

type User struct {
	Id       uint   `db:"id"`
	Username string `db:"username"`
	Avatar   string `db:"avatar"`
}

func (u *User) Name_() string {
	return dbschema.WithPrefix(`nging_user`)
}

type Customer struct {
	Id     uint64 `db:"id"`
	Name   string `db:"name"`
	Avatar string `db:"avatar"`
}

func (c *Customer) Name_() string {
	return dbschema.WithPrefix(`official_customer`)
}
