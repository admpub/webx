package meilisearch

import (
	"fmt"

	"github.com/admpub/webx/application/library/search"
	"github.com/meilisearch/meilisearch-go"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// https://docs.meilisearch.com/learn/getting_started/quick_start.html#download-and-launch

func New(cfg search.Config) (search.Searcher, error) {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:    cfg.Host,
		APIKey:  cfg.Password,
		Timeout: cfg.Timeout,
	})
	return &MeiliSearch{
		Client: client,
	}, nil
}

type MeiliSearch struct {
	*meilisearch.Client
}

func (m *MeiliSearch) getIndex(index string) *meilisearch.Index {
	indexInstance := m.Client.Index(index)
	return indexInstance
}

func (m *MeiliSearch) Add(index string, primaryKey string, docs ...interface{}) error {
	indexM := m.getIndex(index)
	var (
		t   *meilisearch.TaskInfo
		err error
	)
	if len(primaryKey) > 0 {
		t, err = indexM.AddDocuments(docs, primaryKey)
	} else {
		t, err = indexM.AddDocuments(docs)
	}
	if err != nil {
		return err
	}
	if !isOk(t.Status) {
		err = fmt.Errorf(`%+v`, *t)
	}
	return err
}

func (m *MeiliSearch) Update(index string, primaryKey string, docs ...interface{}) error {
	indexM := m.getIndex(index)
	var (
		t   *meilisearch.TaskInfo
		err error
	)
	if len(primaryKey) > 0 {
		t, err = indexM.UpdateDocuments(docs, primaryKey)
	} else {
		t, err = indexM.UpdateDocuments(docs)
	}
	if err != nil {
		return err
	}
	if !isOk(t.Status) {
		err = fmt.Errorf(`%+v`, *t)
	}
	return err
}

func (m *MeiliSearch) GetTask(index string, taskUID int64) error {
	indexM := m.getIndex(index)
	task, err := indexM.GetTask(taskUID)
	com.Dump(task)
	return err
}

func (m *MeiliSearch) Delete(index string, ids ...string) error {
	indexM := m.getIndex(index)
	_, err := indexM.DeleteDocuments(ids)
	return err
}

func (m *MeiliSearch) Flush() error {
	return nil
}

func (m *MeiliSearch) DeleteIndex(index string) error {
	_, err := m.Client.DeleteIndex(index)
	return err
}

func (m *MeiliSearch) InitIndex(cfg *search.IndexConfig) error {
	// m.DeleteIndex(cfg.Index)
	if _, err := m.getIndex(cfg.Index).FetchInfo(); err != nil {
		m.Client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        cfg.Index,
			PrimaryKey: cfg.PrimaryKey,
		})
		searchableAttributes := cfg.SearchableAttributes
		sortableAttributes := cfg.SortableAttributes
		filterableAttributes := cfg.FilterableAttributes
		index := m.getIndex(cfg.Index)
		if len(searchableAttributes) > 0 {
			index.UpdateSearchableAttributes(&searchableAttributes)
		}
		if len(sortableAttributes) > 0 {
			index.UpdateSortableAttributes(&sortableAttributes)
		}
		if len(filterableAttributes) > 0 {
			index.UpdateFilterableAttributes(&filterableAttributes)
		}
	}
	return nil
}

func (m *MeiliSearch) Search(index string, keywords string, options *search.SearchRequest) (int64, []echo.H, error) {
	searchRes, err := m.getIndex(index).Search(
		keywords,
		&meilisearch.SearchRequest{
			Offset: options.Offset,
			Limit:  options.Limit,
			Filter: options.Filter,
			Sort:   options.SortFields,
		},
	)
	if err != nil {
		return 0, nil, err
	}

	rows := make([]echo.H, len(searchRes.Hits))
	for i, v := range searchRes.Hits {
		vm := echo.H(v.(map[string]interface{}))
		m := echo.H{`Doc`: vm}
		if vm.Has(`_formatted`) {
			m.Set(`_formatted`, vm.GetStore(`_formatted`))
		}
		rows[i] = m
	}
	return searchRes.EstimatedTotalHits, rows, nil
}
