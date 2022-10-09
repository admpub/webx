package schemas

type Index struct {
	Name        string         `json:"name"`
	StorageType string         `json:"storage_type"`
	Mappings    *IndexMappings `json:"mappings"`
}

type IndexMappings struct {
	Properties *IndexProperty `json:"properties"`
}

type IndexProperty map[string]*IndexPropertyT

type IndexPropertyT struct {
	Type           string `json:"type"`
	Index          bool   `json:"index"`
	Store          bool   `json:"store"`
	Sortable       bool   `json:"sortable"`
	Aggregatable   bool   `json:"aggregatable"`
	Highlightable  bool   `json:"highlightable"`
	Analyzer       string `json:"analyzer"`
	SearchAnalyzer string `json:"search_analyzer"`
	Format         string `json:"format"`
}

type IndexListItem struct {
	Name        string                 `json:"name"`
	StorageType string                 `json:"storage_type"`
	Mappings    *IndexMappings         `json:"mappings"`
	Settings    map[string]interface{} `json:"settings"`
	CreateAt    string                 `json:"create_at"`    // "0001-01-01T00:00:00Z",
	UpdateAt    string                 `json:"update_at"`    // "0001-01-01T00:00:00Z",
	DocsCount   uint                   `json:"docs_count"`   // 36935,
	StorageSize uint                   `json:"storage_size"` // 11301
}

type IndexList []*IndexListItem
