package advert

import "github.com/admpub/webx/application/dbschema"

type ItemAndPosition struct {
	*dbschema.OfficialAdItem
	Rendered   string
	AdPosition *dbschema.OfficialAdPosition `db:"-,relation=id:position_id|gtZero"`
}

func (i *ItemAndPosition) GetWidth() uint {
	return i.AdPosition.Width
}

func (i *ItemAndPosition) GetHeight() uint {
	return i.AdPosition.Height
}

func (i *ItemAndPosition) GetURL() string {
	return i.Url
}

func (i *ItemAndPosition) GetContent() string {
	return i.Content
}

func (i *ItemAndPosition) GetContype() string {
	return i.Contype
}

type ItemAndRendered struct {
	*dbschema.OfficialAdItem
	Rendered string
}

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
	i.Rendered = Render(i)
	return i
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
