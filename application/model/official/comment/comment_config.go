package comment

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/webx/application/dbschema"
)

type CommentConfig struct {
	SN2ID       func(ctx echo.Context, sn string) (uint64, error) // 序列号转ID
	GetTarget   func(ctx echo.Context, targetID uint64) (res echo.H, breadcrumb []echo.KV, detailURL string, err error)
	CheckMaster func(echo.Context, *dbschema.OfficialCommonComment) error
	AfterAdd    func(echo.Context, *dbschema.OfficialCommonComment) error
	AfterEdit   func(echo.Context, *dbschema.OfficialCommonComment) error
	AfterDelete func(echo.Context, *dbschema.OfficialCommonComment) error
	WithTarget  func(
		ctx echo.Context,
		listx []*CommentAndExtra,
		productIdOwnerIds map[string]map[string]map[string][]uint64, //{source_table:{source_id:{user:[1,2,3]}}}
		targets map[string][]uint64,
		targetObject map[uint64][]int,
	) ([]*CommentAndExtra, error)
	Ident string
	Label string
}

var (
	//CommentAllowUsers all-所有人;buyer-当前商品买家;author-当前文章作者;admin-管理员;allAgent-所有代理;curAgent-当前产品代理;none-无人;designated-指定人员
	CommentAllowUsers = []echo.KV{
		{K: `all`, V: `所有人`},
		{K: `buyer`, V: `当前商品买家`},
		{K: `author`, V: `当前文章作者`},
		{K: `curAgent`, V: `当前产品代理商`},
		{K: `allAgent`, V: `所有代理商`},
		{K: `admin`, V: `管理员`},
		{K: `none`, V: `无`},
		//{K: `designated`, V: `指定人员`},
	}
	//CommentAllowTypes 允许的类型标识
	CommentAllowTypes = map[string]*CommentConfig{
		`article`: {
			GetTarget:   commentArticleGetTarget,
			CheckMaster: commentArticleCheck,
			AfterAdd:    commentArticleAfterAdd,
			AfterDelete: commentArticleAfterDelete,
			WithTarget:  commentArticleWithTarget,
			Ident:       `article`,
			Label:       `资讯`,
		},
		`other`: {
			Ident: `other`,
			Label: `其它`,
		},
	}
)

func CommentSetting() string {
	return config.Setting().GetStore(`base`).String(`comment`)
}

func CommentOpen() bool {
	return CommentSetting() != `close`
}

func CommentReview() bool {
	if !CommentOpen() {
		return true
	}
	return CommentSetting() == `review`
}
