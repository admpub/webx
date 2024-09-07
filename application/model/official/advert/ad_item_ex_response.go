package advert

import (
	"html/template"

	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/com"
)

type ItemResponse struct {
	Content  string `json:"content" xml:"content"`
	Contype  string `json:"contype" xml:"contype"`
	URL      string `json:"url" xml:"url"`
	Start    uint   `json:"start,omitempty" xml:"start,omitempty"`
	End      uint   `json:"end,omitempty" xml:"end,omitempty"`
	Width    uint   `json:"width,omitempty" xml:"width,omitempty"`
	Height   uint   `json:"height,omitempty" xml:"height,omitempty"`
	Rendered string `json:"rendered,omitempty" xml:"rendered,omitempty"`
}

func (i *ItemResponse) GetWidth() uint {
	return i.Width
}

func (i *ItemResponse) GetHeight() uint {
	return i.Height
}

func (i *ItemResponse) GetURL() string {
	return i.URL
}

func (i *ItemResponse) GetContent() string {
	return i.Content
}

func (i *ItemResponse) GetContype() string {
	return i.Contype
}

func (i *ItemResponse) GenHTML() *ItemResponse {
	if i == nil {
		return i
	}
	i.Rendered = Render(i)
	return i
}

func (i *ItemResponse) HTML() template.HTML {
	if i == nil {
		return template.HTML(``)
	}
	if len(i.Rendered) == 0 {
		i.Rendered = Render(i)
	}
	return template.HTML(i.Rendered)
}

func NewItemResponse(item *dbschema.OfficialAdItem, position *dbschema.OfficialAdPosition) *ItemResponse {
	return &ItemResponse{
		Content: item.Content,
		Contype: item.Contype,
		URL:     item.Url,
		Start:   item.Start,
		End:     item.End,
		Width:   position.Width,
		Height:  position.Height,
	}
}

type ItemsResponse []*ItemResponse

func (i ItemsResponse) Rand() *ItemResponse {
	if len(i) < 1 {
		return nil
	}
	return com.RandSlicex(i)
}

func (i ItemsResponse) Valid() bool {
	return len(i) > 0
}

func (i ItemsResponse) Size() int {
	return len(i)
}

func (i ItemsResponse) First() *ItemResponse {
	if len(i) < 1 {
		return nil
	}
	return i[0]
}

func (i ItemsResponse) Last() *ItemResponse {
	if len(i) < 1 {
		return nil
	}
	return i[len(i)-1]
}

func (i *ItemsResponse) GenHTML() *ItemsResponse {
	for _, ad := range *i {
		ad.GenHTML()
	}
	return i
}

func (c *ItemsResponse) Place(_ string) *ItemsResponse {
	return c
}

type PositionAdverts map[string]ItemsResponse

func (c *PositionAdverts) GenHTML() *PositionAdverts {
	for _, adList := range *c {
		adList.GenHTML()
	}
	return c
}

func (c PositionAdverts) Place(ident string) ItemsResponse {
	return c[ident]
}
