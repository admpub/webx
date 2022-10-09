package doc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/admpub/go-zinc/doc/schemas"
	resty "github.com/admpub/resty/v2"
)

type ZincDocSDK interface {
	CreateIndex(name string, p *schemas.IndexProperty, storageType ...string) error
	ListIndex() (schemas.IndexList, error)
	ExistIndex(name string) (bool, error)
	BulkPush(docs []map[string]interface{}) error
	InsertDocumentWithID(name string, id string, doc interface{}) error
	InsertDocument(index string, doc interface{}) error
	DeleteDocument(index string, id string) error
	UpdateDocument(index string, id string, doc interface{}) error
	SearchDocuments(index string, req *schemas.SearchRequest) (*schemas.SearchResponse, error)
}

type zincDocImpl struct {
	client *resty.Client
	host   string
}

func NewSDK(host, user, pwd string, timeout ...time.Duration) (ZincDocSDK, error) {
	client := resty.New()
	client.SetBasicAuth(user, pwd)
	client.SetBaseURL(host)
	client.SetDisableWarn(true)
	if len(timeout) > 0 {
		client.SetTimeout(timeout[0])
	}
	return &zincDocImpl{
		client: client,
		host:   host,
	}, nil
}

func (sdk *zincDocImpl) request() *resty.Request {
	return sdk.client.R()
}

func (c *zincDocImpl) CreateIndex(name string, p *schemas.IndexProperty, storageType ...string) error {
	data := &schemas.Index{
		Name:        name,
		StorageType: "disk",
		Mappings: &schemas.IndexMappings{
			Properties: p,
		},
	}
	if len(storageType) > 0 && len(storageType[0]) > 0 {
		data.StorageType = storageType[0]
	}
	resp, err := c.request().SetBody(data).Put("/api/index")
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("code=%d, msg=%s", resp.StatusCode(), string(resp.Body()))
	}
	return err
}

func (c *zincDocImpl) ListIndex() (schemas.IndexList, error) {
	list := schemas.IndexList{}
	resp, err := c.request().SetResult(&list).Get("/api/index")
	if err != nil {
		return list, err
	}
	if !resp.IsSuccess() {
		return list, fmt.Errorf("code=%d, msg=%s", resp.StatusCode(), string(resp.Body()))
	}
	return list, nil
}

func (c *zincDocImpl) ExistIndex(name string) (bool, error) {
	list, err := c.ListIndex()
	if err != nil {
		return false, err
	}
	for _, item := range list {
		if item.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func (c *zincDocImpl) InsertDocumentWithID(name string, id string, doc interface{}) error {
	resp, err := c.request().SetBody(doc).Put(fmt.Sprintf("/api/%s/_doc/%s", name, id))
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("code=%d, msg=%s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

// documention: https://docs.zincsearch.com/api/document/bulk/
func (c *zincDocImpl) BulkPush(docs []map[string]interface{}) error {
	var buf bytes.Buffer
	for _, doc := range docs {
		b, err := json.Marshal(doc)
		if err == nil {
			buf.Write(b)
			buf.WriteString("\n")
		}
	}
	resp, err := c.request().SetBody(buf.String()).Post("/api/_bulk")
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("code=%d, msg=%s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

func (sdk *zincDocImpl) InsertDocument(index string, doc interface{}) error {
	resp, err := sdk.request().SetBody(doc).Put(fmt.Sprintf("/api/%s/_doc", index))
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("code=%d, msg=%s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

func (sdk *zincDocImpl) DeleteDocument(index string, id string) error {
	resp, err := sdk.request().Delete(fmt.Sprintf("/api/%s/_doc/%s", index, id))
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("code=%d, msg=%s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

func (sdk *zincDocImpl) UpdateDocument(index string, id string, doc interface{}) error {
	resp, err := sdk.request().SetBody(doc).Post(fmt.Sprintf("/api/%s/_update/%s", index, id))
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("code=%d, msg=%s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

func (sdk *zincDocImpl) SearchDocuments(index string, req *schemas.SearchRequest) (*schemas.SearchResponse, error) {
	out := &schemas.SearchResponse{}
	resp, err := sdk.request().SetBody(req).SetResult(out).Post(fmt.Sprintf("/api/%s/_search", index))
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("code=%d, msg=%s", resp.StatusCode(), string(resp.Body()))
	}
	return out, nil
}
