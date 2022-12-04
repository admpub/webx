package doc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/admpub/go-zinc/doc/schemas"
	resty "github.com/admpub/resty/v2"
	"github.com/webx-top/restyclient"
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
	path   string
}

func NewSDK(host, user, pwd string, timeout ...time.Duration) (ZincDocSDK, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	host = u.Host
	if len(u.Scheme) > 0 {
		host = u.Scheme + `://` + host
	} else {
		host = `https://` + host
	}
	path := u.Path
	client := resty.New()
	client.SetBasicAuth(user, pwd)
	client.SetBaseURL(host)
	client.SetDisableWarn(true)
	if len(timeout) > 0 {
		client.SetTimeout(timeout[0])
	}
	path = strings.TrimRight(path, `/`)
	restyclient.InitRestyHook(client)
	return &zincDocImpl{
		client: client,
		host:   host,
		path:   path,
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
	resp, err := c.request().SetBody(data).Put(c.path + "/api/index")
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
	resp, err := c.request().SetResult(&list).Get(c.path + "/api/index")
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
	resp, err := c.request().SetBody(doc).Put(c.path + fmt.Sprintf("/api/%s/_doc/%s", name, id))
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
	resp, err := c.request().SetBody(buf.String()).Post(c.path + "/api/_bulk")
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("code=%d, msg=%s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

func (c *zincDocImpl) InsertDocument(index string, doc interface{}) error {
	resp, err := c.request().SetBody(doc).Put(c.path + fmt.Sprintf("/api/%s/_doc", index))
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("code=%d, msg=%s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

func (c *zincDocImpl) DeleteDocument(index string, id string) error {
	resp, err := c.request().Delete(c.path + fmt.Sprintf("/api/%s/_doc/%s", index, id))
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("code=%d, msg=%s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

func (c *zincDocImpl) UpdateDocument(index string, id string, doc interface{}) error {
	resp, err := c.request().SetBody(doc).Post(c.path + fmt.Sprintf("/api/%s/_update/%s", index, id))
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("code=%d, msg=%s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

func (c *zincDocImpl) SearchDocuments(index string, req *schemas.SearchRequest) (*schemas.SearchResponse, error) {
	out := &schemas.SearchResponse{}
	resp, err := c.request().SetBody(req).SetResult(out).Post(c.path + fmt.Sprintf("/api/%s/_search", index))
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("code=%d, msg=%s", resp.StatusCode(), string(resp.Body()))
	}
	return out, nil
}
