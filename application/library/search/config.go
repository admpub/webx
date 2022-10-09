package search

import "time"

type Config struct {
	Host     string
	User     string
	Password string
	Timeout  time.Duration
}

type IndexConfig struct {
	Index                string
	PrimaryKey           string
	SearchableAttributes []string
	SortableAttributes   []string
	FilterableAttributes []string
	Properties           interface{}
}

type SearchRequest struct {
	SearchType   string
	Offset       int64
	Limit        int64
	Filter       interface{}
	StartTime    time.Time
	EndTime      time.Time
	SortFields   []string
	SearchFields []string
}
