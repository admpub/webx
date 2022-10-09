package official

import (
	"net/url"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/webx/application/dbschema"
)

var (
	FriendlinkMaxUnsuccessfuls int64 = 5 // 每个用户提交友情链接的最大数量（不包含已经通过的）
)

func NewFriendlink(ctx echo.Context) *Friendlink {
	m := &Friendlink{
		OfficialCommonFriendlink: dbschema.NewOfficialCommonFriendlink(ctx),
	}
	return m
}

type Friendlink struct {
	*dbschema.OfficialCommonFriendlink
}

func (f *Friendlink) Exists(customerID uint64, url string, host string) (bool, error) {
	cond := db.NewCompounds()
	if customerID > 0 {
		cond.Add(db.Cond{`customer_id`: customerID})
	}
	cond.Add(db.Or(
		db.Cond{`url`: url},
		db.Cond{`host`: host},
	))
	return f.OfficialCommonFriendlink.Exists(nil, cond.And())
}

func (f *Friendlink) ListPage(cond *db.Compounds, orderby ...interface{}) ([]*FriendlinkExt, error) {
	list := []*FriendlinkExt{}
	_, err := common.NewLister(f, &list, func(r db.Result) db.Result {
		return r.OrderBy(orderby...)
	}, cond.And()).Paging(f.Context())
	return list, err
}

// ListShowAndRecord 前台列表显示，并记录回调
func (f *Friendlink) ListShowAndRecord(limit int, categoryIds ...uint) ([]*dbschema.OfficialCommonFriendlink, error) {
	if !f.Context().Internal().Bool(`FriendlinkReturnRecorded`) {
		f.RecordReturn()
		f.Context().Internal().Set(`FriendlinkReturnRecord`, true)
	}
	if limit == 0 {
		return nil, nil
	}
	list := []*dbschema.OfficialCommonFriendlink{}
	cond := db.NewCompounds()
	cond.AddKV(`process`, `success`)
	var sorts []interface{}
	if len(categoryIds) > 0 && categoryIds[0] > 0 {
		cond.AddKV(`category_id`, categoryIds[0])
	} else {
		sorts = append(sorts, `category_id`)
	}
	sorts = append(sorts, `-return_count`)
	sorts = append(sorts, `-id`)
	_, err := f.ListByOffset(&list, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, 0, limit, cond.And())
	return list, err
}

func (f *Friendlink) ExistsOther(customerID uint64, url string, host string, id uint) (bool, error) {
	cond := db.NewCompounds()
	if customerID > 0 {
		cond.Add(db.Cond{`customer_id`: customerID})
	}
	cond.Add(db.Or(
		db.Cond{`url`: url},
		db.Cond{`host`: host},
	))
	cond.Add(db.Cond{`id`: db.NotEq(id)})
	return f.OfficialCommonFriendlink.Exists(nil, cond.And())
}

func (f *Friendlink) check() error {
	var exists bool
	if len(f.Url) == 0 {
		return f.Context().NewError(code.InvalidParameter, `请输入网址`)
	}
	if !com.IsURL(f.Url) {
		return f.Context().NewError(code.InvalidParameter, `请输入正确的网址`)
	}
	urlInfo, err := url.Parse(f.Url)
	if err != nil {
		return f.Context().NewError(code.InvalidParameter, `网址格式错误: %v`, err.Error())
	}
	if len(urlInfo.Host) == 0 {
		return f.Context().NewError(code.InvalidParameter, `网址格式错误: 域名解析失败`)
	}
	f.Host = urlInfo.Host
	if f.Id > 0 {
		exists, err = f.ExistsOther(f.CustomerId, f.Url, f.Host, f.Id)
	} else {
		if f.CustomerId > 0 {
			count, err := f.Count(nil, db.And(
				db.Cond{`customer_id`: f.CustomerId},
				db.Cond{`process`: db.NotEq(`success`)},
			))
			if err != nil {
				return err
			}
			if count >= FriendlinkMaxUnsuccessfuls {
				return f.Context().NewError(code.DataProcessing, `您已经有“%v”个链接正在处理中，请等待处理通过后再添加新链接`, count)
			}
		}
		exists, err = f.Exists(f.CustomerId, f.Url, f.Host)
	}
	if err != nil {
		return err
	}
	if exists {
		return f.Context().NewError(code.DataAlreadyExists, `网址或域名已经存在`)
	}
	return nil
}

func (f *Friendlink) Add() (pk interface{}, err error) {
	if err := f.check(); err != nil {
		return nil, err
	}
	return f.OfficialCommonFriendlink.Insert()
}

func (f *Friendlink) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.OfficialCommonFriendlink.Update(mw, args...)
}

func (f *Friendlink) IncrReturnCount(id uint) error {
	return f.UpdateFields(nil, echo.H{
		"return_count": db.Raw(`return_count+1`),
		"return_time":  uint(time.Now().Unix()),
	}, `id`, id)
}

func (f *Friendlink) VerifyFail(id uint) error {
	return f.UpdateFields(nil, echo.H{
		"verify_fail_count": db.Raw(`verify_fail_count+1`),
		"verify_time":       uint(time.Now().Unix()),
		"verify_result":     `invalid`,
	}, `id`, id)
}

func (f *Friendlink) RecordReturn() error {
	referer := f.Context().Referer()
	if len(referer) == 0 {
		return nil
	}
	info, err := url.Parse(referer)
	if err != nil {
		return err
	}
	if len(info.Host) == 0 {
		return nil
	}
	if info.Host == f.Context().Host() {
		return nil
	}
	err = f.Get(func(r db.Result) db.Result {
		return r.Select(`host`, `return_time`)
	}, `host`, info.Host)
	if err != nil {
		return err
	}
	if (time.Now().Unix() - 3600) <= int64(f.ReturnTime) { // 每小时只统计一次
		return nil
	}

	return f.IncrReturnCount(f.Id)
}
