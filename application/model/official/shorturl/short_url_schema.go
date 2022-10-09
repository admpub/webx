package shorturl

import (
	"github.com/admpub/webx/application/dbschema"
)

type ShortURLVisitWithURL struct {
	*dbschema.OfficialShortUrlVisit
	Num uint64 `db:"num" bson:"num" comment:"数量" json:"num" xml:"num"`
	URL *dbschema.OfficialShortUrl
}
