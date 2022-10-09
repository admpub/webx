package top

import (
	"github.com/gosimple/slug"
	"github.com/webx-top/com"
)

var SlugifyMaxWidth = 255

func Slugify(v string) string {
	return com.Substr(slug.Make(v), ``, SlugifyMaxWidth)
}
