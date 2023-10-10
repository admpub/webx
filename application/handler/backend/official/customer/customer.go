package customer

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/model/official"
	modelAgent "github.com/admpub/webx/application/model/official/agent"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
	modelLevel "github.com/admpub/webx/application/model/official/level"
)

func Index(ctx echo.Context) error {
	groupId := ctx.Formx(`groupId`).Uint()
	levelId := ctx.Formx(`levelId`).Uint()
	roleId := ctx.Formx(`roleId`).Uint()
	isAgent := ctx.Formx(`isAgent`).Bool()
	online := ctx.Form(`online`)
	m := modelCustomer.NewCustomer(ctx)
	cond := &db.Compounds{}
	if groupId > 0 {
		cond.AddKV(`group_id`, groupId)
	}
	if isAgent {
		cond.AddKV(`agent_level`, db.Gt(0))
	} else {
		agentLevelId := ctx.Formx(`agentLevelId`).Uint()
		if agentLevelId > 0 {
			cond.AddKV(`agent_level`, agentLevelId)
		}
	}
	if levelId > 0 {
		cond.AddKV(`level_id`, levelId)
	}
	common.SelectPageCond(ctx, cond, `id`, `name%`)
	if len(online) > 0 {
		cond.AddKV(`online`, online)
	}
	if roleId > 0 {
		cond.Add(mysql.FindInSet("role_ids", param.AsString(roleId)))
	}
	list := []*modelCustomer.CustomerAndGroup{}
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, &list, func(r db.Result) db.Result {
		return r.Select(factory.DBIGet().OmitSelect(m.OfficialCustomer, `password`, `salt`, `safe_pwd`)...).OrderBy(`-id`).Relation(`Roles`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
			return sel.Columns(`id`, `name`)
		})
	}, cond.And()))
	if err == nil && len(list) > 0 {
		relM := modelLevel.NewRelation(ctx)
		customerIDs := make([]uint64, len(list))
		customerIDIndexes := make(map[uint64]int, len(list))
		for k, v := range list {
			customerIDs[k] = v.Id
			customerIDIndexes[v.Id] = k
		}
		var rels map[uint64][]*modelLevel.RelationExt
		rels, err = relM.ListByCustomerIDs(customerIDs)
		if err == nil {
			for customerID, rows := range rels {
				k := customerIDIndexes[customerID]
				list[k].LevelExtend = rows
			}
		}
	}
	ret := handler.Err(ctx, err)
	ctx.Set(`listData`, list)

	mg := official.NewGroup(ctx)
	var groupList []*dbschema.OfficialCommonGroup
	mg.ListByOffset(&groupList, nil, 0, -1, db.Cond{`type`: `customer`})
	ctx.Set(`groupList`, groupList)

	lv := modelLevel.NewLevel(ctx)
	var levelList []*dbschema.OfficialCustomerLevel
	lv.ListByOffset(&levelList, nil, 0, -1, db.Cond{`group`: `base`})
	ctx.Set(`levelList`, levelList)

	alv := modelAgent.NewAgentLevel(ctx)
	var agentLevelList []*dbschema.OfficialCustomerAgentLevel
	alv.ListByOffset(&agentLevelList, nil, 0, -1)
	ctx.Set(`agentLevelList`, agentLevelList)

	roleM := modelCustomer.NewRole(ctx)
	var roleList []*dbschema.OfficialCustomerRole
	roleM.ListByOffset(&roleList, nil, 0, -1)
	ctx.Set(`roleList`, roleList)

	ctx.Set(`groupId`, groupId)
	ctx.SetFunc(`levelGroupName`, modelLevel.GroupList.Get)
	return ctx.Render(`official/customer/index`, ret)
}

func formFilter(options ...formfilter.Options) echo.FormDataFilter {
	options = append(
		options,
		formfilter.Exclude(`Licenses`, `Created`, `Updated`),
		formfilter.JoinValues(`RoleIds`),
	)
	return formfilter.Build(options...)
}

func Add(ctx echo.Context) error {
	var (
		err error
		id  uint
	)
	m := modelCustomer.NewCustomer(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCustomer, formFilter())
		if err != nil {
			goto END
		}
		password2 := ctx.Form(`password2`)
		safePwd2 := ctx.Form(`safePwd2`)
		if password2 != m.Password {
			err = ctx.E(`两次输入的登录密码不一致`)
			goto END
		}
		if safePwd2 != m.SafePwd {
			err = ctx.E(`两次输入的安全密码不一致`)
			goto END
		}
		if len(ctx.FormValues(`roleIds`)) == 0 {
			m.RoleIds = ``
		}
		_, err = m.Add()
		if err != nil {
			goto END
		}
		handler.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(handler.URLFor(`/official/customer/index`))
	}
	id = ctx.Formx(`copyId`).Uint()
	if id > 0 {
		err = m.Get(nil, `id`, id)
		if err != nil {
			goto END
		}
		echo.StructToForm(ctx, m.OfficialCustomer, ``, echo.LowerCaseFirstLetter)
		ctx.Request().Form().Set(`id`, `0`)
	}

END:
	ctx.Set(`activeURL`, `/official/customer/index`)
	setFormData(ctx, m)
	ctx.Set(`isEdit`, false)
	return ctx.Render(`official/customer/edit`, handler.Err(ctx, err))
}

func setFormData(ctx echo.Context, m *modelCustomer.Customer) {
	mg := official.NewGroup(ctx)
	var groupList []*dbschema.OfficialCommonGroup
	mg.ListByOffset(&groupList, nil, 0, -1, db.Cond{`type`: `customer`})
	ctx.Set(`groupList`, groupList)

	roleM := modelCustomer.NewRole(ctx)
	roleM.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`id`, `name`, `description`)
	}, 0, -1, db.And(db.Cond{`parent_id`: 0}))
	ctx.Set(`roleList`, roleM.Objects())

	var roleIds []uint
	if len(m.RoleIds) > 0 {
		roleIds = param.StringSlice(strings.Split(m.RoleIds, `,`)).Uint()
	}
	ctx.SetFunc(`isChecked`, func(roleId uint) bool {
		for _, rid := range roleIds {
			if rid == roleId {
				return true
			}
		}
		return false
	})
}

func Edit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint64()
	m := modelCustomer.NewCustomer(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.OfficialCustomer, formFilter(formfilter.Exclude(`Name`, `Salt`)))
		modifyPwd := ctx.Form(`modifyPwd`)
		if modifyPwd == `1` {
			password2 := ctx.Form(`password2`)
			if password2 != m.Password {
				err = ctx.E(`两次输入的登录密码不一致`)
			}
		}
		if err == nil {
			modifySafePwd := ctx.Form(`modifySafePwd`)
			if modifySafePwd == `1` {
				safePwd2 := ctx.Form(`safePwd2`)
				if safePwd2 != m.SafePwd {
					err = ctx.E(`两次输入的安全密码不一致`)
				}
			}
		}
		if err == nil {
			m.Id = id
			if len(ctx.FormValues(`roleIds`)) == 0 {
				m.RoleIds = ``
			}
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/official/customer/index`))
			}
		}
		ctx.Request().Form().Set(`name`, m.Name)
	} else if err == nil {
		echo.StructToForm(ctx, m.OfficialCustomer, ``, func(topName, fieldName string) string {
			switch fieldName {
			case `Password`, `SafePwd`:
				return ``
			}
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
	}

	ctx.Set(`activeURL`, `/official/customer/index`)
	setFormData(ctx, m)
	ctx.Set(`isEdit`, true)
	return ctx.Render(`official/customer/edit`, handler.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint64()
	m := modelCustomer.NewCustomer(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/official/customer/index`))
}

// Kick 踢🦶客户下线
func Kick(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint64()
	m := modelCustomer.NewCustomer(ctx)
	err := m.Get(func(r db.Result) db.Result {
		return r.Select(`session_id`)
	}, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if len(m.SessionId) == 0 {
		handler.SendFail(ctx, ctx.T(`此客户没有 session id 记录`))
	} else {
		deviceM := modelCustomer.NewDevice(ctx)
		err = deviceM.Kick(id)
		if err == nil {
			ctx.Session().RemoveID(m.SessionId)
			m.UpdateField(nil, `session_id`, ``, `id`, id)
			handler.SendOk(ctx, ctx.T(`操作成功`))
		} else {
			handler.SendFail(ctx, err.Error())
		}
	}

	return ctx.Redirect(handler.URLFor(`/official/customer/index`))
}
