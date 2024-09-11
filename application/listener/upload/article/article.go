package article

import (
	"fmt"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"

	"github.com/coscms/webcore/library/fileupdater/listener"
	"github.com/coscms/webcore/registry/upload/thumb"
	"github.com/admpub/webx/application/dbschema"
	"github.com/admpub/webx/application/model/official"
	modelArticle "github.com/admpub/webx/application/model/official/article"
)

func init() {
	// - official_common_article
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.OfficialCommonArticle)
		tableID = fmt.Sprint(fm.Id)
		content = fm.ImageOriginal
		property = listener.NewPropertyWith(
			fm,
			db.Cond{`id`: fm.Id},
			listener.FieldValueWith(`image`, thumb.DefaultSize.ThumbValue()),
		)
		return
	}, false).SetTable(`official_common_article`, `image_original`, `image`).ListenDefault()

	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.OfficialCommonArticle)
		tableID = fmt.Sprint(fm.Id)
		content = fm.Content
		property = listener.NewPropertyWith(fm, db.Cond{`id`: fm.Id})
		return
	}, true).SetTable(`official_common_article`, `content`).ListenDefault()

	dbschema.DBI.On(`deleting`, func(m factory.Model, _ ...string) error {
		fm := m.(*dbschema.OfficialCommonArticle)
		var err error
		if len(fm.Tags) > 0 {
			tagsM := official.NewTags(m.Context())
			err = tagsM.DecrNum(modelArticle.GroupName, strings.Split(fm.Tags, `,`))
		}
		return err
	}, `official_common_article`)
}
