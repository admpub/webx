package article

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/echo"
	stdCode "github.com/webx-top/echo/code"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/library/top"
	"github.com/admpub/webx/application/middleware/sessdata"
	modelAuthor "github.com/admpub/webx/application/model/author"
	"github.com/admpub/webx/application/model/official"
	modelArticle "github.com/admpub/webx/application/model/official/article"
	modelComment "github.com/admpub/webx/application/model/official/comment"
	modelCustomer "github.com/admpub/webx/application/model/official/customer"
)

func Detail(c echo.Context) error {
	var err error
	id := c.Paramx(`id`).Uint64()
	op := c.Param(`op`)
	commentURLLayout := c.Echo().URI(`article.detailWithOp`, echo.H{
		`id`: id,
		`op`: `comments`,
	}) + `?page={page}&size={size}&rows={rows}`

	//commentURLLayout := `/article/comment_list?html=1&_pjax=true&page={page}&size={size}&rows={rows}&id=` + param.AsString(id)
	flat := true
	c.Set(`flat`, flat)
	c.SetFunc(`commentList`, func() []*modelComment.CommentAndExtra {
		c.Request().Form().Set(`_pjax`, `true`)
		c.Request().Form().Set(`pageSize`, `10`)
		commentList, _ := articleCommentList(c, id, ``, `article`, ``, flat, commentURLLayout, `Comment`)
		//c.Request().Form().Del(`_pjax`)
		return commentList
	})
	if op == `comments` && (c.IsPjax() || c.IsAjax()) {
		data := c.Data()
		b, err := c.Fetch(`/article/comment_list`, nil)
		if err != nil {
			return err
		}
		data.SetData(echo.H{
			`html`:       string(b),
			`pagination`: c.Get(`pagination`),
		})
		return c.JSON(data)
	}

	// 资讯
	articleM := modelArticle.NewArticle(c)
	err = articleM.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return c.NewError(stdCode.DataNotFound, `不存id为“%d”的资讯`, id)
		}
		return err
	}
	customer := sessdata.Customer(c)
	if articleM.Display == `N` && (customer == nil || articleM.OwnerType != `customer` || customer.Id != articleM.OwnerId) {
		return c.NewError(stdCode.DataUnavailable, `此文章不可查看`)
	}
	articleM.Content = top.HideContent(articleM.Content, articleM.Contype, modelArticle.GetContentHideDetector(customer, articleM.OfficialCommonArticle))
	c.Set(`data`, articleM.OfficialCommonArticle)
	categories, err := articleM.GetCategories()
	if err != nil {
		return err
	}
	c.Set(`categories`, categories)
	var (
		sourceInfo echo.KV
		listURL    = sessdata.URLFor(`/articles`)
	)
	if len(articleM.SourceTable) > 0 && len(articleM.SourceId) > 0 {
		infoGetter := modelArticle.Source.GetInfoGetter(articleM.SourceTable)
		if infoGetter != nil {
			sourceInfo, err = infoGetter(c, articleM.SourceId)
			if err != nil {
				return err
			}
			listURL = sessdata.URLFor(`/articlesBy/` + articleM.SourceTable + `/` + articleM.SourceId)
			c.Set(`sourceTable`, articleM.SourceTable)
			c.Set(`sourceId`, articleM.SourceId)
		}
	}
	c.Set(`sourceInfo`, sourceInfo)
	c.Set(`needReviewComment`, articleM.OfficialCommonArticle.CommentAutoDisplay == `N`)
	c.Set(`targetSubtype`, `article`)
	c.Set(`targetType`, ``)
	author := modelAuthor.New(articleM.OwnerId, articleM.OwnerType).Get(c)
	c.Set(`author`, author)
	err = articleM.UpdateField(nil, `views`, db.Raw(`views+1`), `id`, id)
	if err != nil {
		log.Error(err)
		err = nil
	}
	tmpl := `detail`
	if len(articleM.Template) > 0 {
		tmpl = articleM.Template
	}
	c.SetFunc(`relationList`, articleM.RelationList)
	c.SetFunc(`queryList`, articleM.QueryList)
	c.SetFunc(`tagList`, func() []*dbschema.OfficialCommonTags {
		tags, _ := getTags(c)
		return tags
	})

	// 资讯点赞记录
	clickFlowM := official.NewClickFlow(c)
	var (
		ownerID   uint64
		ownerType string
	)
	if customer != nil {
		ownerID = customer.Id
		ownerType = `customer`
	}
	clickFlowM.Find(`article`, articleM.Id, ownerID, ownerType)
	c.Set(`clickFlow`, clickFlowM.OfficialCommonClickFlow)

	// 是否允许评论
	c.Set(`disabledCommentMessage`, articleM.IsAllowedComment(customer))
	extraCond := db.NewCompounds()
	extraCond.Add(db.Cond{`source_id`: articleM.SourceId})
	extraCond.Add(db.Cond{`source_table`: articleM.SourceTable})
	nextRow, _ := articleM.NextRow(articleM.Id, extraCond)
	c.Set(`nextRow`, nextRow)
	prevRow, _ := articleM.PrevRow(articleM.Id, extraCond)
	c.Set(`prevRow`, prevRow)
	c.Set(`listURL`, listURL+c.DefaultExtension())
	return c.Render(`article/`+tmpl, handler.Err(c, err))
}

func ArticleListBy(c echo.Context) error {
	sourceID := c.Param(`sourceId`)
	sourceTable := c.Param(`sourceTable`)
	categoryID := c.Formx(`categoryId`).Uint()
	return ListBy(c, sourceID, sourceTable, categoryID)
}

func ListBy(c echo.Context, sourceID string, sourceTable string, categoryID ...uint) error {
	tag := c.Query(`tag`)
	articleM := modelArticle.NewArticle(c)
	articleM.SourceId = sourceID
	articleM.SourceTable = sourceTable
	cond := db.NewCompounds()
	cond.Add(db.Cond{`source_id`: sourceID})
	cond.Add(db.Cond{`source_table`: sourceTable})
	cond.Add(db.Cond{`display`: `Y`})
	var categories []dbschema.OfficialCommonCategory
	if len(categoryID) > 0 && categoryID[0] > 0 {
		cond.Add(db.Or(
			db.Cond{`category1`: categoryID[0]},
			db.Cond{`category2`: categoryID[0]},
			db.Cond{`category3`: categoryID[0]},
			db.Cond{`category_id`: categoryID[0]},
		))
		cateM := official.NewCategory(c)
		categories, _ = cateM.Parents(categoryID[0])
	}
	if len(tag) > 0 {
		cond.Add(articleM.TagCond(tag))
	}
	c.Request().Form().Set(`pageSize`, `20`)
	articles, err := articleM.ListPageSimple(cond)
	if err != nil {
		return err
	}
	c.Set(`articles`, articles)
	c.Set(`categories`, categories)
	var sourceInfo echo.KV
	if len(sourceTable) > 0 && len(sourceID) > 0 {
		infoGetter := modelArticle.Source.GetInfoGetter(sourceTable)
		if infoGetter != nil {
			sourceInfo, err = infoGetter(c, sourceID)
			if err != nil {
				return err
			}
		}
	}
	c.Set(`sourceInfo`, sourceInfo)
	c.Set(`sourceTable`, sourceTable)
	c.Set(`sourceId`, sourceID)
	c.Set(`listURL`, sessdata.URLFor(`/articlesBy/`+sourceTable+`/`+sourceID))
	c.SetFunc(`relationList`, articleM.RelationList)
	c.SetFunc(`queryList`, articleM.QueryList)
	c.SetFunc(`tagList`, func() []*dbschema.OfficialCommonTags {
		tags, _ := getTags(c)
		return tags
	})
	return c.Render(`article/list_by`, handler.Err(c, err))
}

func List(c echo.Context) error {
	var err error
	tag := c.Query(`tag`)
	query := c.Form(`q`)
	articleM := modelArticle.NewArticle(c)
	cond := db.NewCompounds()
	cond.Add(db.Cond{`display`: `Y`})
	categoryID := c.Queryx(`categoryId`).Uint()
	var categories []dbschema.OfficialCommonCategory
	if categoryID > 0 {
		cond.Add(db.Or(
			db.Cond{`category1`: categoryID},
			db.Cond{`category2`: categoryID},
			db.Cond{`category3`: categoryID},
			db.Cond{`category_id`: categoryID},
		))
		cateM := official.NewCategory(c)
		categories, _ = cateM.Parents(categoryID)
	}
	if len(tag) > 0 {
		cond.Add(articleM.TagCond(tag))
	}
	if len(query) > 0 {
		cond.From(mysql.SearchField(`title`, query))
	}
	c.Request().Form().Set(`pageSize`, `20`)
	articles, err := articleM.ListPageSimple(cond)
	if err != nil {
		return err
	}
	c.Set(`articles`, articles)
	c.Set(`categories`, categories)
	c.Set(`tag`, tag)
	c.Set(`listURL`, sessdata.URLFor(`/articles`)+c.DefaultExtension())
	c.SetFunc(`relationList`, articleM.RelationList)
	c.SetFunc(`queryList`, articleM.QueryList)
	c.SetFunc(`tagList`, func() []*dbschema.OfficialCommonTags {
		tags, _ := getTags(c)
		return tags
	})
	return c.Render(`article/list`, handler.Err(c, err))
}

func Pay(c echo.Context) error {
	customer := sessdata.Customer(c)
	if customer == nil {
		return common.ErrUserNotLoggedIn
	}
	id := c.Paramx(`id`).Uint64()

	// 资讯
	articleM := modelArticle.NewArticle(c)
	err := articleM.Get(func(r db.Result) db.Result {
		return r.Select(`id`, `price`, `title`)
	}, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return c.NewError(stdCode.DataNotFound, `不存id为“%d”的资讯`, id)
		}
		return err
	}
	if articleM.Price <= 0 {
		return c.NewError(stdCode.Failure, `此为免费文章，无需购买`)
	}
	walletM := modelCustomer.NewWallet(c)
	exists, err := walletM.Flow.Exists(nil, db.And(
		db.Cond{`customer_id`: customer.Id},
		db.Cond{`asset_type`: modelCustomer.AssetTypeMoney},
		db.Cond{`amount_type`: modelCustomer.AmountTypeBalance},
		db.Cond{`source_type`: `buy`},
		db.Cond{`source_table`: `official_common_article`},
		db.Cond{`source_id`: articleM.Id},
	))
	if err != nil {
		return err
	}
	if exists {
		return c.NewError(stdCode.RepeatOperation, c.T(`您已经支付过了，请不要重复购买`))
	}
	walletM.Flow.CustomerId = customer.Id
	walletM.Flow.AssetType = modelCustomer.AssetTypeMoney
	walletM.Flow.AmountType = modelCustomer.AmountTypeBalance
	walletM.Flow.Amount = -articleM.Price
	walletM.Flow.SourceType = `buy`
	walletM.Flow.SourceTable = `official_common_article`
	walletM.Flow.SourceId = articleM.Id
	walletM.Flow.TradeNo = ``
	walletM.Flow.Status = modelCustomer.FlowStatusConfirmed //状态(pending-待确认;confirmed-已确认;canceled-已取消)
	walletM.Flow.Description = `购买文章: ` + articleM.Title
	err = walletM.AddFlow()
	return c.JSON(c.Data().SetError(err))
}
