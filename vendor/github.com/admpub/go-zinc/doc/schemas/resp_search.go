package schemas

import "time"

type SearchResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Hits     struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Index     string                 `json:"_index"`
			Type      string                 `json:"_type"`
			ID        string                 `json:"_id"`
			Score     float64                `json:"_score"`
			Timestamp time.Time              `json:"@timestamp"`
			Source    map[string]interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
	Error string `json:"error"`
}
