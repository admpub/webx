package schemas

import "time"

type SearchRequest struct {
	SearchType string `json:"search_type"`
	Query      struct {
		Term      string    `json:"term"`
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
	} `json:"query"`
	SortFields []string `json:"sort_fields"`
	From       int      `json:"from"`
	MaxResults int      `json:"max_results"`
	Source     []string `json:"_source"`
}
