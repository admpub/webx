package shorturl

import (
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/webx/application/dbschema"
	modelShorturl "github.com/admpub/webx/application/model/official/shorturl"
)

func applyFormData(ctx echo.Context, murl *dbschema.OfficialShortUrl) error {
	if murl.Id == 0 {
		murl.OwnerType = `user`
	}
	expired := ctx.Formx(`expired`).String()
	murl.Expired = 0
	if len(expired) > 0 {
		if t, err := time.Parse(`2006-01-02`, expired); err == nil && !t.IsZero() {
			murl.Expired = uint(t.Unix())
		}
	}
	murl.Available = ctx.Formx(`available`).String()
	murl.LongUrl = ctx.Formx(`longUrl`).String()
	murl.Password = ctx.Formx(`password`).String()
	return nil
}

// Index 短网址列表
func Index(ctx echo.Context) error {
	//user := handler.User(ctx)
	var err error
	m := modelShorturl.NewShortURL(ctx)
	cond := db.Compounds{
		//db.Cond{`owner_id`: user.Id},
		//db.Cond{`owner_type`: `user`},
	}
	shortID := ctx.Formx(`q`).String()
	if len(shortID) > 0 {
		cond.AddKV(`short_url`, shortID)
	}
	sorts := common.Sorts(ctx, `official_short_url`, `-id`)
	_, err = common.NewLister(m.URL, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(ctx)
	ctx.Set(`listData`, m.URL.Objects())
	return ctx.Render(`official/short_url/index`, common.Err(ctx, err))
}

// Add 创建短网址
func Add(ctx echo.Context) error {
	user := handler.User(ctx)
	var (
		err error
		id  uint64
	)
	m := modelShorturl.NewShortURL(ctx)
	if ctx.IsPost() {
		m.URL.OwnerId = uint64(user.Id)
		if err = applyFormData(ctx, m.URL); err != nil {
			goto END
		}
		_, err = m.Add()
		if err != nil {
			goto END
		}
		common.SendOk(ctx, ctx.T(`添加成功`))
		return ctx.Redirect(handler.URLFor(`/official/short_url/index`))
	}
	id = ctx.Formx(`copyId`).Uint64()
	if id > 0 {
		err = m.URL.Get(nil, `id`, id)
		if err == nil {
			echo.StructToForm(ctx, m.URL, ``, echo.LowerCaseFirstLetter)
			ctx.Request().Form().Set(`id`, `0`)
			if m.URL.Expired > 0 {
				ctx.Request().Form().Set(`expired`, time.Unix(int64(m.URL.Expired), 0).Format(`2006-01-02`))
			}
		}
	}

END:
	ctx.Set(`activeURL`, `/official/short_url/index`)
	ctx.Set(`title`, ctx.T(`添加短链接`))
	return ctx.Render(`official/short_url/edit`, common.Err(ctx, err))
}

// Edit 修改短网址
func Edit(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint64()
	if id < 1 {
		return ctx.E(`参数“%s”值无效`, `id`)
	}
	m := modelShorturl.NewShortURL(ctx)
	err := m.URL.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.E(`短网址不存在`)
		}
		return err
	}
	if ctx.IsPost() {
		if err = applyFormData(ctx, m.URL); err != nil {
			goto END
		}
		err = m.Edit(nil, `id`, m.URL.Id)
		if err != nil {
			goto END
		}
		common.SendOk(ctx, ctx.T(`修改成功`))
		return ctx.Redirect(handler.URLFor(`/official/short_url/index`))
	}
	echo.StructToForm(ctx, m.URL, ``, echo.LowerCaseFirstLetter)
	if m.URL.Expired > 0 {
		ctx.Request().Form().Set(`expired`, time.Unix(int64(m.URL.Expired), 0).Format(`2006-01-02`))
	}

END:
	ctx.Set(`activeURL`, `/official/short_url/index`)
	ctx.Set(`title`, ctx.T(`修改短链接`))
	return ctx.Render(`official/short_url/edit`, common.Err(ctx, err))
}

// Delete 删除短网址
func Delete(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint64()
	if id < 1 {
		return ctx.E(`参数“%s”值无效`, `id`)
	}
	m := modelShorturl.NewShortURL(ctx)
	err := m.URL.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.E(`短网址不存在`)
		}
		return err
	}
	err = m.URL.Delete(nil, `id`, id)
	if err != nil {
		return err
	}
	common.SendOk(ctx, ctx.T(`删除成功`))
	return ctx.Redirect(handler.URLFor(`/official/short_url/index`))
}

// Analysis 用户短网址访问统计
func Analysis(ctx echo.Context) error {
	m := modelShorturl.NewShortURL(ctx)

	// - lastList

	cond := db.Compounds{}
	_, err := m.Visit.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`-created`)
	}, 0, 50, cond.And())
	if err != nil {
		return err
	}
	rows := m.Visit.Objects()
	lastList, err := m.VisitListWithURL(rows)
	if err != nil {
		return err
	}
	ctx.Set(`lastList`, lastList)

	// - isp top10
	var ispTopList []*modelShorturl.ShortURLVisitWithURL
	_, err = m.Visit.ListByOffset(&ispTopList, func(r db.Result) db.Result {
		return r.Select(db.Raw(`isp,COUNT(isp) AS num`)).OrderBy(`-num`).Group(`isp`)
	}, 0, 10, cond.And())
	if err != nil {
		return err
	}
	ispTopList, err = m.VisitListFillData(ispTopList)
	if err != nil {
		return err
	}
	ctx.Set(`ispTopList`, ispTopList)

	// - province top10

	var provinceTopList []*modelShorturl.ShortURLVisitWithURL
	_, err = m.Visit.ListByOffset(&provinceTopList, func(r db.Result) db.Result {
		return r.Select(db.Raw(`province,COUNT(province) AS num`)).OrderBy(`-num`).Group(`province`)
	}, 0, 10, cond.And())
	if err != nil {
		return err
	}
	ctx.Set(`provinceTopList`, provinceTopList)

	// - browser top10

	var browserTopList []*modelShorturl.ShortURLVisitWithURL
	_, err = m.Visit.ListByOffset(&browserTopList, func(r db.Result) db.Result {
		return r.Select(db.Raw(`browser,COUNT(browser) AS num`)).OrderBy(`-num`).Group(`browser`)
	}, 0, 10, cond.And())
	if err != nil {
		return err
	}
	ctx.Set(`browserTopList`, browserTopList)

	// - os top10

	var osTopList []*modelShorturl.ShortURLVisitWithURL
	_, err = m.Visit.ListByOffset(&osTopList, func(r db.Result) db.Result {
		return r.Select(db.Raw(`os,COUNT(os) AS num`)).OrderBy(`-num`).Group(`os`)
	}, 0, 10, cond.And())
	if err != nil {
		return err
	}
	ctx.Set(`osTopList`, osTopList)

	//ctx.Set(`activeURL`, `/official/short_url/analysis`)
	return ctx.Render(`official/short_url/analysis`, common.Err(ctx, err))
}
