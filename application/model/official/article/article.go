package article

import (
	"fmt"
	"strings"
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/download/downloadByContent"
	"github.com/admpub/webx/application/library/top"
	"github.com/admpub/webx/application/library/xcommon"
	"github.com/admpub/webx/application/library/xrole"
	"github.com/admpub/webx/application/library/xrole/xroleutils"
	"github.com/admpub/webx/application/model/official"
	"github.com/admpub/webx/application/model/official/agent"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

func NewArticle(ctx echo.Context) *Article {
	m := &Article{
		OfficialCommonArticle: dbschema.NewOfficialCommonArticle(ctx),
	}
	return m
}

type Article struct {
	*dbschema.OfficialCommonArticle
	DisallowCreateTags bool
}

func (f *Article) ListByTag(recv interface{}, mw func(db.Result) db.Result, page, size int, tag string) (func() int64, error) {
	return f.OfficialCommonArticle.List(recv, mw, page, size, f.TagCond(tag))
}

func (f *Article) GetCategories(rows ...*dbschema.OfficialCommonArticle) ([]dbschema.OfficialCommonCategory, error) {
	row := f.OfficialCommonArticle
	if len(rows) > 0 {
		row = rows[0]
	}
	cateM := official.NewCategory(f.Context())
	return cateM.Parents(row.CategoryId)
}

func (f *Article) TagCond(tag string) db.Compound {
	return official.TagCond(tag)
}

func (f *Article) check(old *dbschema.OfficialCommonArticle) error {
	if len(f.Title) < 1 {
		return f.Context().NewError(code.InvalidParameter, `标题不能为空`).SetZone(`title`)
	}
	if f.Price < 0 {
		return f.Context().NewError(code.InvalidParameter, `价格不能为负数`).SetZone(`price`)
	}
	if f.CategoryId > 0 {
		cateM := official.NewCategory(f.Context())
		parentIDs := cateM.ParentIds(f.CategoryId)
		f.Category3 = 0
		f.Category2 = 0
		f.Category1 = 0
		switch len(parentIDs) {
		case 4:
			//f.CategoryId = parentIDs[3]
			fallthrough
		case 3:
			f.Category3 = parentIDs[2]
			fallthrough
		case 2:
			f.Category2 = parentIDs[1]
			fallthrough
		case 1:
			f.Category1 = parentIDs[0]
		}
	}
	var err error
	var oldTags []string
	if old != nil && len(old.Tags) > 0 {
		oldTags = strings.Split(old.Tags, `,`)
	}
	tagsM := official.NewTags(f.Context())
	tagsM.Use(f.Trans())
	tags, err := tagsM.UpdateTags(f.Id == 0, GroupName, oldTags, strings.Split(f.Tags, `,`), f.DisallowCreateTags)
	if err != nil {
		return err
	}
	f.Tags = strings.Join(tags, `,`)

	if len(f.Contype) == 0 || !Contype.Has(f.Contype) {
		f.Contype = `text`
	}
	f.Content = common.ContentEncode(f.Content, f.Contype)
	f.Slugify = top.Slugify(f.Title)
	return err
}

func (f *Article) SyncRemoteImage() (string, error) {
	newContent, err := downloadByContent.SyncRemoteImage(f.Context(), `default`, fmt.Sprint(f.Id), f.Content, f.Contype)
	return newContent, err
}

func (f *Article) CustomerTodayCount(customerID interface{}) (int64, error) {
	startTs, endTs := top.TodayTimestamp()
	return f.Count(nil, db.And(
		db.Cond{`owner_type`: `customer`},
		db.Cond{`owner_id`: customerID},
		db.Cond{`created`: db.Between(startTs, endTs)},
	))
}

func (f *Article) CustomerPendingCount(customerID interface{}) (int64, error) {
	return f.Count(nil, db.And(
		db.Cond{`owner_type`: `customer`},
		db.Cond{`owner_id`: customerID},
		db.Cond{`display`: `N`},
	))
}

func (f *Article) CustomerPendingTodayCount(customerID interface{}) (int64, error) {
	startTs, endTs := top.TodayTimestamp()
	return f.Count(nil, db.And(
		db.Cond{`owner_type`: `customer`},
		db.Cond{`owner_id`: customerID},
		db.Cond{`display`: `N`},
		db.Cond{`created`: db.Between(startTs, endTs)},
	))
}

func (f *Article) checkCustomerAdd(permission *xrole.RolePermission) error {
	err := xcommon.CheckRoleCustomerAdd(f.Context(), permission, BehaviorName, f.OwnerId, f)
	if err == nil {
		return err
	}
	switch err {
	case xcommon.ErrCustomerAddClosed:
		return f.Context().E(`文章投稿功能已关闭`)
	case xcommon.ErrCustomerAddMaxPerDay:
		return f.Context().E(`投稿失败。您的账号已达到今日最大投稿数量`)
	case xcommon.ErrCustomerAddMaxPending:
		return f.Context().E(`投稿失败。您的待审核文章数量已达上限，请等待审核通过后再投稿`)
	default:
		return err
	}
}

func (f *Article) Add() (pk interface{}, err error) {
	if f.OwnerType == `customer` {
		permission := xroleutils.CustomerPermission(f.Context())
		if err = f.checkCustomerAdd(permission); err != nil {
			return nil, err
		}
	}
	f.Context().Begin()
	if err = f.check(nil); err != nil {
		f.Context().Rollback()
		return
	}
	syncRemoteImage := f.Context().Formx(`syncRemoteImage`).Bool()
	if syncRemoteImage {
		f.Content, err = f.SyncRemoteImage()
		if err != nil {
			f.Context().Rollback()
			return
		}
	}
	pk, err = f.OfficialCommonArticle.Insert()
	if err != nil {
		f.Context().Rollback()
		return
	}
	err = f.Context().Commit()
	return
}

func (f *Article) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	old := dbschema.NewOfficialCommonArticle(f.Context())
	err := old.Get(nil, args...)
	if err != nil {
		return err
	}
	if err := f.check(old); err != nil {
		return err
	}
	syncRemoteImage := f.Context().Formx(`syncRemoteImage`).Bool()
	if syncRemoteImage {
		f.Content, err = f.SyncRemoteImage()
		if err != nil {
			return err
		}
	}
	err = f.OfficialCommonArticle.Update(mw, args...)
	return err
}

func (f *Article) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	err := f.OfficialCommonArticle.Delete(mw, args...)
	return err
}

var OrderQuerier = func(c echo.Context, customer *dbschema.OfficialCustomer, sourceId, sourceTable string) error {
	return c.E(`很抱歉，您不是当前商品买家，本文只有当前商品买家才能评论`)
}

// IsAllowedComment 是否可以评论
func (f *Article) IsAllowedComment(customer *dbschema.OfficialCustomer) error {
	articleM := f
	c := f.Context()
	var err error
	if articleM.CloseComment == `Y` {
		return c.E(`本文已经关闭评论`)
	}
	switch articleM.CommentAllowUser {
	case `all`:
	case `buyer`:
		if customer == nil {
			return c.E(`很抱歉，您不是当前商品买家，本文只有当前商品买家才能评论`)
		}
		boughtDetector := Source.GetBoughtDetector(articleM.SourceTable)
		if boughtDetector != nil {
			if boughtDetector(customer, articleM.SourceId) {
				return nil
			}
			return c.E(`很抱歉，您不是当前商品买家，本文只有当前商品买家才能评论`)
		}
		return OrderQuerier(c, customer, articleM.SourceId, articleM.SourceTable)
	case `author`:
		user := handler.User(c)
		if articleM.OwnerType == `user` {
			if user == nil || articleM.OwnerId != uint64(user.Id) {
				customerM := modelCustomer.NewCustomer(c)
				err = customerM.Get(nil, db.And(
					db.Cond{`id`: customer.Id},
					db.Cond{`uid`: articleM.OwnerId},
				))
				if err != nil {
					if err == db.ErrNoMoreRows {
						return c.E(`很抱歉，您不是本文作者，本文只有作者才能评论`)
					}
					return err
				}
			}
		} else { //客户发布的文章，后台用户的作者可以评论
			if user == nil && (customer == nil || customer.Id != articleM.OwnerId) {
				return c.E(`很抱歉，您不是本文作者，本文只有作者才能评论`)
			}
		}
	case `curAgent`:
		if customer == nil {
			return c.E(`很抱歉，您没有代理当前商品，本文只有当前商品的代理商才能评论`)
		}
		agentDetector := Source.GetAgentDetector(articleM.SourceTable)
		if agentDetector != nil {
			if agentDetector(customer, articleM.SourceId) {
				return nil
			}
			return c.E(`很抱歉，您没有代理当前商品，本文只有当前商品的代理商才能评论`)
		}
		agentProdM := agent.NewAgentProduct(c)
		err = agentProdM.Get(nil, db.And(
			db.Cond{`agent_id`: customer.Id},
			db.Cond{`product_id`: articleM.SourceId},
			db.Cond{`product_table`: articleM.SourceTable},
			db.Or(
				db.Cond{`expired`: 0},
				db.Cond{`expired`: db.Gt(time.Now().Unix())},
			),
		))
		if err != nil {
			if err == db.ErrNoMoreRows {
				return c.E(`很抱歉，您没有代理当前商品，本文只有当前商品的代理商才能评论`)
			}
			return err
		}
	case `allAgent`:
		if customer == nil {
			return c.E(`很抱歉，您不是代理商，本文只有代理商才能评论`)
		}
		agentDetector := Source.GetAgentDetector(articleM.SourceTable)
		if agentDetector != nil {
			if agentDetector(customer, ``) {
				return nil
			}
			return c.E(`很抱歉，您不是代理商，本文只有代理商才能评论`)
		}
		agentProdM := agent.NewAgentProduct(c)
		err = agentProdM.Get(nil, db.And(
			db.Cond{`agent_id`: customer.Id},
			db.Or(
				db.Cond{`expired`: 0},
				db.Cond{`expired`: db.Gt(time.Now().Unix())},
			),
		))
		if err != nil {
			if err == db.ErrNoMoreRows {
				return c.E(`很抱歉，您不是代理商，本文只有代理商才能评论`)
			}
			return err
		}
	case `admin`:
		user := handler.User(c)
		if user == nil || user.Id != 1 {
			return c.E(`很抱歉，本文只有管理员才能评论`)
		}
	case `none`:
		return c.E(`很抱歉，本文暂未开放评论`)
	}
	return err
}

func (f *Article) ListPageSimple(cond *db.Compounds, orderby ...interface{}) ([]*ArticleWithOwner, error) {
	if len(orderby) == 0 {
		orderby = []interface{}{`-id`}
	}
	rows := []*ArticleWithOwner{}
	_, err := common.NewLister(f, &rows, func(r db.Result) db.Result {
		return r.Select(factory.DBIGet().OmitSelect(f.OfficialCommonArticle, `content`)...).OrderBy(orderby...)
	}, cond.And()).Paging(f.Context())
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (f *Article) ListByOffsetSimple(cond *db.Compounds, limit int, offset int, orderby ...interface{}) ([]*ArticleWithOwner, error) {
	if len(orderby) == 0 {
		orderby = []interface{}{`-id`}
	}
	rows := []*ArticleWithOwner{}
	_, err := f.OfficialCommonArticle.ListByOffset(&rows, func(r db.Result) db.Result {
		return r.Select(factory.DBIGet().OmitSelect(f.OfficialCommonArticle, `content`)...).OrderBy(orderby...)
	}, offset, limit, cond.And())
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (f *Article) ListPage(cond *db.Compounds, orderby ...interface{}) ([]*ArticleAndSourceInfo, error) {
	list := []*ArticleAndSourceInfo{}
	_, err := common.NewLister(f, &list, func(r db.Result) db.Result {
		return r.Select(factory.DBIGet().OmitSelect(f, `content`)...).OrderBy(orderby...)
	}, cond.And()).Paging(f.Context())
	if err != nil {
		return list, err
	}
	err = WithSourceInfo(f.Context(), list)
	if err != nil {
		return nil, err
	}
	tg := make([]official.ICategory, len(list))
	for k, v := range list {
		tg[k] = v
	}
	categoryM := official.NewCategory(f.Context())
	err = categoryM.FillTo(tg)
	return list, err
}

func (f *Article) NextRow(currentID uint64, extraCond *db.Compounds) (*dbschema.OfficialCommonArticle, error) {
	row := dbschema.NewOfficialCommonArticle(nil)
	row.CPAFrom(f.OfficialCommonArticle)
	cond := db.NewCompounds()
	cond.AddKV(`display`, `Y`)
	cond.AddKV(`id`, db.Lt(currentID))
	if extraCond != nil {
		cond.From(extraCond)
	}
	err := row.Get(func(r db.Result) db.Result {
		return r.Select(`id`, `title`, `image`, `created`).OrderBy(`-id`)
	}, cond.And())
	return row, err
}

func (f *Article) PrevRow(currentID uint64, extraCond *db.Compounds) (*dbschema.OfficialCommonArticle, error) {
	row := dbschema.NewOfficialCommonArticle(nil)
	row.CPAFrom(f.OfficialCommonArticle)
	cond := db.NewCompounds()
	cond.AddKV(`display`, `Y`)
	cond.AddKV(`id`, db.Gt(currentID))
	if extraCond != nil {
		cond.From(extraCond)
	}
	err := row.Get(func(r db.Result) db.Result {
		return r.Select(`id`, `title`, `image`, `created`).OrderBy(`id`)
	}, cond.And())
	return row, err
}

func (f *Article) RelationList(limit int, orderby ...interface{}) []*ArticleWithOwner {
	cond := db.NewCompounds()
	if len(f.SourceTable) > 0 {
		cond.Add(db.Cond{`source_table`: f.SourceTable})
		if len(f.SourceId) > 0 {
			cond.Add(db.Cond{`source_id`: f.SourceId})
		}
	}
	return f.CommonQueryList(cond, limit, 0, orderby...)
}

func (f *Article) CommonQueryList(cond *db.Compounds, limit int, offset int, orderby ...interface{}) []*ArticleWithOwner {
	cond.Add(db.Cond{`display`: `Y`})
	rows, _ := f.ListByOffsetSimple(cond, limit, offset, orderby...)
	return rows
}

func (f *Article) QueryList(query string, limit int, offset int, orderby ...interface{}) []*ArticleWithOwner {
	cond := db.NewCompounds()
	r := common.NewSortedURLValues(query)
	r.ApplyCond(cond)
	return f.CommonQueryList(cond, limit, offset, orderby...)
}
