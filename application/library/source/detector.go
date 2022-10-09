package source

import (
	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/echo"
)

type (
	// Detector 确认逻辑 如果sourceId为空白字符串，则代表不限制具体资源
	// 例如 BoughtDetector(customer, ``) 则代表判断是否购买过任一产品
	Detector func(customer *dbschema.OfficialCustomer, sourceId string) bool
)

func (s *Source) GetBoughtDetector(sourceTable string) Detector {
	item := s.GetItem(sourceTable)
	if item == nil {
		return nil
	}
	switch v := item.X.(type) {
	case *Info:
		return v.BoughtDetector()
	case Info:
		return (&v).BoughtDetector()
	default:
		return nil
	}
}

func (s *Source) GetAgentDetector(sourceTable string) Detector {
	item := s.GetItem(sourceTable)
	if item == nil {
		return nil
	}
	switch v := item.X.(type) {
	case *Info:
		return v.AgentDetector()
	case Info:
		return (&v).AgentDetector()
	default:
		return nil
	}
}

func (s *Source) GetInfoGetter(sourceTable string) InfoGetter {
	item := s.GetItem(sourceTable)
	if item == nil {
		return nil
	}
	switch v := item.X.(type) {
	case *Info:
		return v.InfoGetter()
	case Info:
		return (&v).InfoGetter()
	default:
		return nil
	}
}

func (s *Source) GetInfoMapGetter(sourceTable string) InfoMapGetter {
	item := s.GetItem(sourceTable)
	if item == nil {
		return nil
	}
	switch v := item.X.(type) {
	case *Info:
		return v.InfoMapGetter()
	case Info:
		return (&v).InfoMapGetter()
	default:
		return nil
	}
}

func (s *Source) GetTagsGetter(sourceTable string) TagsGetter {
	item := s.GetItem(sourceTable)
	if item == nil {
		return nil
	}
	switch v := item.X.(type) {
	case *Info:
		return v.TagsGetter()
	case Info:
		return (&v).TagsGetter()
	default:
		return nil
	}
}

func (s *Source) GetSelectPageHandler(sourceTable string) func(echo.Context) error {
	item := s.GetItem(sourceTable)
	if item == nil {
		return nil
	}
	switch v := item.X.(type) {
	case *Info:
		return v.SelectPage()
	case Info:
		return (&v).SelectPage()
	default:
		return nil
	}
}
