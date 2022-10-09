package source

import "github.com/webx-top/echo"

type (
	// InfoGetter 查询单个数据
	InfoGetter func(ctx echo.Context, sourceId string) (echo.KV, error)
	// InfoMapGetter 查询多个数据
	InfoMapGetter func(ctx echo.Context, sourceId ...string) (map[string]echo.KV, error)
	// TagsGetter 查询标签
	TagsGetter func(ctx echo.Context, sourceId interface{}) ([]echo.H, error)
	// Info 资源信息
	Info struct {
		isBought          Detector
		isAgent           Detector
		getInfo           InfoGetter
		getInfoMap        InfoMapGetter
		getTags           TagsGetter
		selectPageHandler func(echo.Context) error
	}
	// Infor 数据接口
	Infor interface {
		BoughtDetector() Detector
		AgentDetector() Detector
		InfoGetter() InfoGetter
		TagsGetter() TagsGetter
		SelectPage() func(echo.Context) error
	}
)

func NewInfo() *Info {
	return &Info{}
}

func (s *Info) BoughtDetector() Detector {
	return s.isBought
}

func (s *Info) SetBoughtDetector(fn Detector) *Info {
	s.isBought = fn
	return s
}

func (s *Info) AgentDetector() Detector {
	return s.isAgent
}

func (s *Info) SetAgentDetector(fn Detector) *Info {
	s.isAgent = fn
	return s
}

func (s *Info) InfoGetter() InfoGetter {
	return s.getInfo
}

func (s *Info) SetInfoGetter(fn InfoGetter) *Info {
	s.getInfo = fn
	return s
}

func (s *Info) InfoMapGetter() InfoMapGetter {
	return s.getInfoMap
}

func (s *Info) SetInfoMapGetter(fn InfoMapGetter) *Info {
	s.getInfoMap = fn
	return s
}

func (s *Info) TagsGetter() TagsGetter {
	return s.getTags
}

func (s *Info) SetTagsGetter(fn TagsGetter) *Info {
	s.getTags = fn
	return s
}

func (s *Info) SetSelectPageHandler(h func(echo.Context) error) *Info {
	s.selectPageHandler = h
	return s
}

func (s *Info) SelectPage() func(echo.Context) error {
	return s.selectPageHandler
}
