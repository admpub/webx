package search

import "github.com/webx-top/echo"

type Searcher interface {
	Add(index string, primaryKey string, docs ...interface{}) error
	Update(index string, primaryKey string, docs ...interface{}) error
	Delete(index string, ids ...string) error
	Flush() error
	InitIndex(cfg *IndexConfig) error
	Search(index string, keywords string, options *SearchRequest) (int64, []echo.H, error)
}

var DefaultSearch = &NopSearch{}

type NopSearch struct{}

func (n *NopSearch) Add(index string, primaryKey string, docs ...interface{}) error {
	return nil
}

func (m *NopSearch) Update(index string, primaryKey string, docs ...interface{}) error {
	return nil
}

func (m *NopSearch) Delete(index string, ids ...string) error {
	return nil
}

func (m *NopSearch) InitIndex(cfg *IndexConfig) error {
	return nil
}

func (n *NopSearch) Flush() error {
	return nil
}

func (n *NopSearch) Search(index string, keywords string, options *SearchRequest) (int64, []echo.H, error) {
	return 0, nil, nil
}
