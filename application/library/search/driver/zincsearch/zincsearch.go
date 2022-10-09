package zincsearch

import (
	"github.com/admpub/go-zinc/doc"
	"github.com/admpub/go-zinc/doc/schemas"
	"github.com/admpub/webx/application/library/search"
	"github.com/webx-top/echo/param"
)

func New(cfg search.Config) *ZincSearch {
	sdk, err := doc.NewSDK(cfg.Host, cfg.User, cfg.Password, cfg.Timeout)
	if err != nil {
		panic(err)
	}
	return &ZincSearch{
		ZincDocSDK: sdk,
	}
}

type ZincSearch struct {
	doc.ZincDocSDK
}

func (m *ZincSearch) Add(index string, primaryKey string, docs ...interface{}) error {
	var err error
	for _, doc := range docs {
		row := param.AsStore(doc)
		id := row.String(primaryKey)
		err = m.ZincDocSDK.InsertDocumentWithID(index, id, doc)
		if err != nil {
			return err
		}
	}
	return err
}

func (m *ZincSearch) Update(index string, primaryKey string, doc interface{}) error {
	row := param.AsStore(doc)
	id := row.String(primaryKey)
	return m.ZincDocSDK.UpdateDocument(index, id, doc)
}

func (m *ZincSearch) Delete(index string, ids ...string) error {
	var err error
	for _, id := range ids {
		err = m.ZincDocSDK.DeleteDocument(index, id)
		if err != nil {
			return err
		}
	}
	return err
}

func (m *ZincSearch) Flush() error {
	return nil
}

func (m *ZincSearch) InitIndex(cfg *search.IndexConfig) error {
	property := cfg.Properties.(*schemas.IndexProperty)
	return m.ZincDocSDK.CreateIndex(cfg.Index, property)
}

func (m *ZincSearch) Search(index string, keywords string, options *search.SearchRequest) (interface{}, error) {
	cfg := &schemas.SearchRequest{
		SearchType: options.SearchType,
		SortFields: options.SortFields,
		From:       int(options.Offset),
		MaxResults: int(options.Limit),
		Source:     options.SearchFields,
	}
	if len(cfg.SearchType) == 0 {
		cfg.SearchType = `querystring`
	}
	cfg.Query.Term = keywords
	cfg.Query.StartTime = options.StartTime
	cfg.Query.EndTime = options.EndTime
	searchRes, err := m.ZincDocSDK.SearchDocuments(index, cfg)
	if err != nil {
		return nil, err
	}

	return searchRes.Hits, nil
}
