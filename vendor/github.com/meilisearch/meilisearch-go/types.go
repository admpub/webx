package meilisearch

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
)

//
// Internal types to Meilisearch
//

// Client is a structure that give you the power for interacting with an high-level api with Meilisearch.
type Client struct {
	config     ClientConfig
	httpClient *fasthttp.Client
}

// Index is the type that represent an index in Meilisearch
type Index struct {
	UID        string    `json:"uid"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	PrimaryKey string    `json:"primaryKey,omitempty"`
	client     *Client
}

// Return of multiple indexes is wrap in a IndexesResults
type IndexesResults struct {
	Results []Index `json:"results"`
	Offset  int64   `json:"offset"`
	Limit   int64   `json:"limit"`
	Total   int64   `json:"total"`
}

type IndexesQuery struct {
	Limit  int64 `json:"limit,omitempty"`
	Offset int64 `json:"offset,omitempty"`
}

// Settings is the type that represents the settings in Meilisearch
type Settings struct {
	RankingRules         []string            `json:"rankingRules,omitempty"`
	DistinctAttribute    *string             `json:"distinctAttribute,omitempty"`
	SearchableAttributes []string            `json:"searchableAttributes,omitempty"`
	DisplayedAttributes  []string            `json:"displayedAttributes,omitempty"`
	StopWords            []string            `json:"stopWords,omitempty"`
	Synonyms             map[string][]string `json:"synonyms,omitempty"`
	FilterableAttributes []string            `json:"filterableAttributes,omitempty"`
	SortableAttributes   []string            `json:"sortableAttributes,omitempty"`
	TypoTolerance        *TypoTolerance      `json:"typoTolerance,omitempty"`
	Pagination           *Pagination         `json:"pagination,omitempty"`
	Faceting             *Faceting           `json:"faceting,omitempty"`
}

// TypoTolerance is the type that represents the typo tolerance setting in Meilisearch
type TypoTolerance struct {
	Enabled             bool                `json:"enabled,omitempty"`
	MinWordSizeForTypos MinWordSizeForTypos `json:"minWordSizeForTypos,omitempty"`
	DisableOnWords      []string            `json:"disableOnWords,omitempty"`
	DisableOnAttributes []string            `json:"disableOnAttributes,omitempty"`
}

// MinWordSizeForTypos is the type that represents the minWordSizeForTypos setting in the typo tolerance setting in Meilisearch
type MinWordSizeForTypos struct {
	OneTypo  int64 `json:"oneTypo,omitempty"`
	TwoTypos int64 `json:"twoTypos,omitempty"`
}

// Pagination is the type that represents the pagination setting in Meilisearch
type Pagination struct {
	MaxTotalHits int64 `json:"maxTotalHits"`
}

// Faceting is the type that represents the faceting setting in Meilisearch
type Faceting struct {
	MaxValuesPerFacet int64 `json:"maxValuesPerFacet"`
}

// Version is the type that represents the versions in Meilisearch
type Version struct {
	CommitSha  string `json:"commitSha"`
	CommitDate string `json:"commitDate"`
	PkgVersion string `json:"pkgVersion"`
}

// StatsIndex is the type that represent the stats of an index in Meilisearch
type StatsIndex struct {
	NumberOfDocuments int64            `json:"numberOfDocuments"`
	IsIndexing        bool             `json:"isIndexing"`
	FieldDistribution map[string]int64 `json:"fieldDistribution"`
}

// Stats is the type that represent all stats
type Stats struct {
	DatabaseSize int64                 `json:"databaseSize"`
	LastUpdate   time.Time             `json:"lastUpdate"`
	Indexes      map[string]StatsIndex `json:"indexes"`
}

// TaskStatus is the status of a task.
type TaskStatus string

const (
	// TaskStatusUnknown is the default TaskStatus, should not exist
	TaskStatusUnknown TaskStatus = "unknown"
	// TaskStatusEnqueued the task request has been received and will be processed soon
	TaskStatusEnqueued TaskStatus = "enqueued"
	// TaskStatusProcessing the task is being processed
	TaskStatusProcessing TaskStatus = "processing"
	// TaskStatusSucceeded the task has been successfully processed
	TaskStatusSucceeded TaskStatus = "succeeded"
	// TaskStatusFailed a failure occurred when processing the task, no changes were made to the database
	TaskStatusFailed TaskStatus = "failed"
)

// Task indicates information about a task resource
//
// Documentation: https://docs.meilisearch.com/learn/advanced/asynchronous_operations.html
type Task struct {
	Status     TaskStatus          `json:"status"`
	UID        int64               `json:"uid,omitempty"`
	TaskUID    int64               `json:"taskUid,omitempty"`
	IndexUID   string              `json:"indexUid"`
	Type       string              `json:"type"`
	Error      meilisearchApiError `json:"error,omitempty"`
	Duration   string              `json:"duration,omitempty"`
	EnqueuedAt time.Time           `json:"enqueuedAt"`
	StartedAt  time.Time           `json:"startedAt,omitempty"`
	FinishedAt time.Time           `json:"finishedAt,omitempty"`
	Details    Details             `json:"details,omitempty"`
}

// TaskInfo indicates information regarding a task returned by an asynchronous method
//
// Documentation: https://docs.meilisearch.com/reference/api/tasks.html#tasks
type TaskInfo struct {
	Status     TaskStatus          `json:"status"`
	TaskUID    int64               `json:"taskUid,omitempty"`
	IndexUID   string              `json:"indexUid"`
	Type       string              `json:"type"`
	Error      meilisearchApiError `json:"error,omitempty"`
	Duration   string              `json:"duration,omitempty"`
	EnqueuedAt time.Time           `json:"enqueuedAt"`
	StartedAt  time.Time           `json:"startedAt,omitempty"`
	FinishedAt time.Time           `json:"finishedAt,omitempty"`
	Details    Details             `json:"details,omitempty"`
}

// TasksQuery is the request body for list documents method
type TasksQuery struct {
	Limit    int64    `json:"limit,omitempty"`
	From     int64    `json:"from,omitempty"`
	IndexUID []string `json:"indexUid,omitempty"`
	Status   []string `json:"status,omitempty"`
	Type     []string `json:"type,omitempty"`
}

type Details struct {
	ReceivedDocuments    int                 `json:"receivedDocuments,omitempty"`
	IndexedDocuments     int                 `json:"indexedDocuments,omitempty"`
	DeletedDocuments     int                 `json:"deletedDocuments,omitempty"`
	PrimaryKey           string              `json:"primaryKey,omitempty"`
	RankingRules         []string            `json:"rankingRules,omitempty"`
	DistinctAttribute    *string             `json:"distinctAttribute,omitempty"`
	SearchableAttributes []string            `json:"searchableAttributes,omitempty"`
	DisplayedAttributes  []string            `json:"displayedAttributes,omitempty"`
	StopWords            []string            `json:"stopWords,omitempty"`
	Synonyms             map[string][]string `json:"synonyms,omitempty"`
	FilterableAttributes []string            `json:"filterableAttributes,omitempty"`
	SortableAttributes   []string            `json:"sortableAttributes,omitempty"`
}

// Return of multiple tasks is wrap in a TaskResult
type TaskResult struct {
	Results []Task `json:"results"`
	Limit   int64  `json:"limit"`
	From    int64  `json:"from"`
	Next    int64  `json:"next"`
}

// Keys allow the user to connect to the Meilisearch instance
//
// Documentation: https://docs.meilisearch.com/learn/advanced/security.html#protecting-a-meilisearch-instance
type Key struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Key         string    `json:"key,omitempty"`
	UID         string    `json:"uid,omitempty"`
	Actions     []string  `json:"actions,omitempty"`
	Indexes     []string  `json:"indexes,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// This structure is used to send the exact ISO-8601 time format managed by Meilisearch
type KeyParsed struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UID         string    `json:"uid,omitempty"`
	Actions     []string  `json:"actions,omitempty"`
	Indexes     []string  `json:"indexes,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	ExpiresAt   *string   `json:"expiresAt"`
}

// This structure is used to update a Key
type KeyUpdate struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Return of multiple keys is wrap in a KeysResults
type KeysResults struct {
	Results []Key `json:"results"`
	Offset  int64 `json:"offset"`
	Limit   int64 `json:"limit"`
	Total   int64 `json:"total"`
}

type KeysQuery struct {
	Limit  int64 `json:"limit,omitempty"`
	Offset int64 `json:"offset,omitempty"`
}

// Information to create a tenant token
//
// ExpiresAt is a time.Time when the key will expire.
// Note that if an ExpiresAt value is included it should be in UTC time.
// ApiKey is the API key parent of the token.
type TenantTokenOptions struct {
	APIKey    string
	ExpiresAt time.Time
}

// Custom Claims structure to create a Tenant Token
type TenantTokenClaims struct {
	APIKeyUID   string      `json:"apiKeyUid"`
	SearchRules interface{} `json:"searchRules"`
	jwt.RegisteredClaims
}

//
// Request/Response
//

// CreateIndexRequest is the request body for create index method
type CreateIndexRequest struct {
	UID        string `json:"uid,omitempty"`
	PrimaryKey string `json:"primaryKey,omitempty"`
}

// SearchRequest is the request url param needed for a search query.
// This struct will be converted to url param before sent.
//
// Documentation: https://docs.meilisearch.com/reference/features/search_parameters.html
type SearchRequest struct {
	Offset                int64
	Limit                 int64
	AttributesToRetrieve  []string
	AttributesToCrop      []string
	CropLength            int64
	CropMarker            string
	AttributesToHighlight []string
	HighlightPreTag       string
	HighlightPostTag      string
	MatchingStrategy      string
	Filter                interface{}
	ShowMatchesPosition   bool
	Facets                []string
	PlaceholderSearch     bool
	Sort                  []string
}

// SearchResponse is the response body for search method
type SearchResponse struct {
	Hits               []interface{} `json:"hits"`
	EstimatedTotalHits int64         `json:"estimatedTotalHits"`
	Offset             int64         `json:"offset"`
	Limit              int64         `json:"limit"`
	ProcessingTimeMs   int64         `json:"processingTimeMs"`
	Query              string        `json:"query"`
	FacetDistribution  interface{}   `json:"facetDistribution,omitempty"`
}

// DocumentQuery is the request body get one documents method
type DocumentQuery struct {
	Fields []string `json:"fields,omitempty"`
}

// DocumentsQuery is the request body for list documents method
type DocumentsQuery struct {
	Offset int64    `json:"offset,omitempty"`
	Limit  int64    `json:"limit,omitempty"`
	Fields []string `json:"fields,omitempty"`
}

type DocumentsResult struct {
	Results []map[string]interface{} `json:"results"`
	Limit   int64                    `json:"limit"`
	Offset  int64                    `json:"offset"`
	Total   int64                    `json:"total"`
}

// RawType is an alias for raw byte[]
type RawType []byte

// Health is the request body for set Meilisearch health
type Health struct {
	Status string `json:"status"`
}

// UpdateIndexRequest is the request body for update Index primary key
type UpdateIndexRequest struct {
	PrimaryKey string `json:"primaryKey"`
}

// Unknown is unknown json type
type Unknown map[string]interface{}

// UnmarshalJSON supports json.Unmarshaler interface
func (b *RawType) UnmarshalJSON(data []byte) error {
	*b = data
	return nil
}

// MarshalJSON supports json.Marshaler interface
func (b RawType) MarshalJSON() ([]byte, error) {
	return b, nil
}
