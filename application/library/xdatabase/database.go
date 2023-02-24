package xdatabase

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
)

func GetAndLock(m factory.Model, cond db.Compound, middlewares ...func(db.Result) db.Result) error {
	ctx := m.Context()
	cid := m.ConnID()
	nmr := m.Namer()
	p := m.NewParam()
	if len(middlewares) > 0 {
		p = p.SetMW(middlewares[0])
	}
	err := p.AddArgs(cond).Select().Amend(func(queryIn string) string {
		return queryIn + ` FOR UPDATE`
	}).One(m)
	m.SetContext(ctx)
	m.SetConnID(cid)
	m.SetNamer(nmr)
	return err
}

func ListAndLock(m factory.Model, recvPtr interface{}, cond db.Compound, offset int, limit int, middlewares ...func(db.Result) db.Result) error {
	p := m.NewParam()
	if len(middlewares) > 0 {
		p = p.SetMW(middlewares[0])
	}
	err := p.AddArgs(cond).SetOffset(offset).SetSize(limit).Select().Amend(func(queryIn string) string {
		return queryIn + ` FOR UPDATE`
	}).All(recvPtr)
	return err
}
