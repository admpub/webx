package meilisearch

import (
	"github.com/admpub/webx/application/library/search"
	"github.com/meilisearch/meilisearch-go"
)

func New(cfg search.Config) *MeiliSearch {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:    cfg.Host,
		APIKey:  cfg.Password,
		Timeout: cfg.Timeout,
	})
	return &MeiliSearch{
		Client: client,
	}
}

type MeiliSearch struct {
	*meilisearch.Client
}

func (m *MeiliSearch) Add(index string, primaryKey string, docs ...interface{}) error {
	indexM := m.Client.Index(index)
	_, err := indexM.AddDocuments(docs, primaryKey)
	return err
}

func (m *MeiliSearch) Update(index string, primaryKey string, docs ...interface{}) error {
	indexM := m.Client.Index(index)
	_, err := indexM.UpdateDocuments(docs, primaryKey)
	return err
}

func (m *MeiliSearch) Delete(index string, ids ...string) error {
	indexM := m.Client.Index(index)
	_, err := indexM.DeleteDocuments(ids)
	return err
}

func (m *MeiliSearch) Flush() error {
	return nil
}

func (m *MeiliSearch) InitIndex(cfg *search.IndexConfig) error {
	if _, err := m.Client.Index(cfg.Index).FetchInfo(); err != nil {
		m.Client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        cfg.Index,
			PrimaryKey: cfg.PrimaryKey,
		})
		searchableAttributes := cfg.SearchableAttributes
		sortableAttributes := cfg.SortableAttributes
		filterableAttributes := cfg.FilterableAttributes
		index := m.Client.Index(cfg.Index)
		index.UpdateSearchableAttributes(&searchableAttributes)
		index.UpdateSortableAttributes(&sortableAttributes)
		index.UpdateFilterableAttributes(&filterableAttributes)
	}
	return nil
}

func (m *MeiliSearch) Search(index string, keywords string, options *search.SearchRequest) (interface{}, error) {
	searchRes, err := m.Client.Index(index).Search(
		keywords,
		&meilisearch.SearchRequest{
			Offset: options.Offset,
			Limit:  options.Limit,
			Filter: options.Filter,
			Sort:   options.SortFields,
		},
	)
	if err != nil {
		return nil, err
	}

	return searchRes.Hits, nil
}
