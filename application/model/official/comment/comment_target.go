package comment

import (
	"fmt"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/subdomains"

	"github.com/admpub/webx/application/model/official/article"
)

func commentArticleWithTarget(
	ctx echo.Context,
	listx []*CommentAndExtra,
	productIdOwnerIds map[string]map[string]map[string][]uint64,
	targets map[string][]uint64,
	targetObject map[uint64][]int,
) ([]*CommentAndExtra, error) {
	m := article.NewArticle(ctx)
	mw := func(r db.Result) db.Result {
		return r.Select(`id`, `owner_id`, `owner_type`, `source_id`, `source_table`, `title`)
	}
	_, err := m.ListByOffset(nil, mw, 0, -1, `id`, db.In(targets[`article`]))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return listx, err
		}
		return listx, nil
	}
	for _, target := range m.Objects() {
		if _, _y := productIdOwnerIds[target.SourceTable]; !_y {
			productIdOwnerIds[target.SourceTable] = map[string]map[string][]uint64{}
		}
		if _, _y := productIdOwnerIds[target.SourceTable][target.SourceId]; !_y {
			productIdOwnerIds[target.SourceTable][target.SourceId] = map[string][]uint64{}
		}
		if _, _y := productIdOwnerIds[target.SourceTable][target.SourceId][target.OwnerType]; !_y {
			productIdOwnerIds[target.SourceTable][target.SourceId][target.OwnerType] = []uint64{}
		}
		for _, k := range targetObject[target.Id] {
			listx[k].Extra[`targetObject`] = echo.H{
				`ownerId`:      target.OwnerId,
				`ownerType`:    target.OwnerType,
				`productId`:    target.SourceId,
				`productTable`: target.SourceTable,
				`detailURL`:    subdomains.Default.URL(`/article/`+fmt.Sprint(target.Id)+ctx.DefaultExtension(), `frontend`),
				`title`:        target.Title,
			}
			if !com.InUint64Slice(listx[k].OwnerId, productIdOwnerIds[target.SourceTable][target.SourceId][target.OwnerType]) {
				productIdOwnerIds[target.SourceTable][target.SourceId][target.OwnerType] = append(productIdOwnerIds[target.SourceTable][target.SourceId][target.OwnerType], listx[k].OwnerId)
			}
		}
	}
	return listx, err
}
