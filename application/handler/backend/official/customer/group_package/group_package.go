package group_package

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/formbuilder"
	"github.com/coscms/webcore/library/nsql"
	"github.com/coscms/webfront/model/i18nm"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	modelLevel "github.com/coscms/webfront/model/official/level"
)

func Index(ctx echo.Context) error {
	group := ctx.Form(`group`)
	m := modelCustomer.NewGroupPackage(ctx)
	cond := db.NewCompounds()
	if len(group) > 0 {
		cond.AddKV(`group`, group)
	}
	nsql.SelectPageCond(ctx, cond, `id`, `title`)
	_, err := common.PagingWithLister(ctx, common.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ret := common.Err(ctx, err)
	list := m.Objects()
	ctx.Set(`listData`, list)
	ctx.Set(`groupList`, modelLevel.GroupList.Slice())
	ctx.SetFunc(`levelGroupName`, modelLevel.GroupList.Get)
	ctx.SetFunc(`timeUnitName`, modelCustomer.GroupPackageTimeUnits.Get)
	return ctx.Render(`official/customer/group_package/index`, ret)
}

func Add(ctx echo.Context) error {
	var err error
	m := modelCustomer.NewGroupPackage(ctx)
	if ctx.IsGet() {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				m.Id = 0
				i18nm.SetModelTranslationsToForm(ctx, m.OfficialCustomerGroupPackage, uint64(id))
			} else {
				m.Sort = 5000
			}
		} else {
			m.Sort = 5000
		}

	}
	form := formbuilder.New(ctx,
		m.OfficialCustomerGroupPackage,
		formbuilder.ConfigFile(`official/customer/group_package/edit`),
		formbuilder.AllowedNames(
			`iconImage`, `iconClass`, `recommend`, `disabled`, `sort`,
			`timeUnit`, `timeDuration`, `price`, `group`, `title`, `description`,
		),
	)
	form.OnPost(func() error {
		_, err := m.Add()
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(ctx, m.OfficialCustomerGroupPackage, uint64(m.Id), i18nm.OptionContentType(`description`, `text`))
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/official/customer/group_package/index`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()
	nameField := form.MultilingualField(config.FromFile().Language.Default, `title`, `title`)
	nameField.AddTag(`required`)

	ctx.Set(`activeURL`, `/official/customer/group_package/index`)
	ctx.Set(`title`, ctx.T(`添加套餐`))
	setFormData(ctx, m)
	return ctx.Render(`official/customer/group_package/edit`, common.Err(ctx, err))
}

func setFormData(ctx echo.Context, m *modelCustomer.GroupPackage) {
	ctx.Set(`groupList`, modelLevel.GroupList.Slice())
	ctx.Set(`timeUnits`, modelCustomer.GroupPackageTimeUnits.Slice())
}

func Edit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	if id == 0 {
		return ctx.NewError(code.InvalidParameter, `参数错误`).SetZone(`id`)
	}
	m := modelCustomer.NewGroupPackage(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if ctx.IsGet() {
		if ctx.IsAjax() {

			disabled := ctx.Query(`disabled`)
			if len(disabled) > 0 {
				if !common.IsBoolFlag(disabled) {
					return ctx.NewError(code.InvalidParameter, ``).SetZone(`disabled`)
				}
				m.Disabled = disabled
				data := ctx.Data()
				err = m.UpdateField(nil, `disabled`, disabled, db.Cond{`id`: id})
				if err != nil {
					data.SetError(err)
					return ctx.JSON(data)
				}
				data.SetInfo(ctx.T(`操作成功`))
				return ctx.JSON(data)
			}

			recommend := ctx.Query(`recommend`)
			if len(recommend) > 0 {
				if !common.IsBoolFlag(recommend) {
					return ctx.NewError(code.InvalidParameter, ``).SetZone(`recommend`)
				}
				m.Recommend = recommend
				data := ctx.Data()
				err = m.UpdateField(nil, `recommend`, recommend, db.Cond{`id`: id})
				if err != nil {
					data.SetError(err)
					return ctx.JSON(data)
				}
				data.SetInfo(ctx.T(`操作成功`))
				return ctx.JSON(data)
			}

		}
		i18nm.SetModelTranslationsToForm(ctx, m.OfficialCustomerGroupPackage, uint64(id))
	}
	form := formbuilder.New(ctx,
		m.OfficialCustomerGroupPackage,
		formbuilder.ConfigFile(`official/customer/group_package/edit`),
		formbuilder.AllowedNames(
			`iconImage`, `iconClass`, `recommend`, `disabled`, `sort`,
			`timeUnit`, `timeDuration`, `price`, `group`, `title`, `description`,
		),
	)
	form.OnPost(func() error {
		err := m.Edit(nil, db.Cond{`id`: id})
		if err != nil {
			return err
		}
		err = i18nm.SaveModelTranslations(ctx, m.OfficialCustomerGroupPackage, uint64(m.Id), i18nm.OptionContentType(`description`, `text`))
		if err != nil {
			return err
		}
		common.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(backend.URLFor(`/official/customer/group_package/index`))
	})
	err = form.RecvSubmission()
	if form.Exited() {
		return form.Error()
	}
	form.Generate()
	nameField := form.MultilingualField(config.FromFile().Language.Default, `title`, `title`)
	nameField.AddTag(`required`)

	ctx.Set(`activeURL`, `/official/customer/group_package/index`)
	ctx.Set(`title`, ctx.T(`编辑套餐`))
	setFormData(ctx, m)
	return ctx.Render(`official/customer/group_package/edit`, common.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := modelCustomer.NewGroupPackage(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/official/customer/group_package/index`))
}
