package meilisearch

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/webx/application/library/search"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

// curl -X GET 'http://localhost:7700/indexes' -H 'Authorization: Bearer masterKey'

// curl -X DELETE 'http://localhost:7700/indexes/default' -H 'Content-Type: application/json' -H 'Authorization: Bearer masterKey'

// curl -X GET 'http://localhost:7700/indexes/default/documents' -H 'Authorization: Bearer masterKey'

/*
curl \
  -X PUT 'http://localhost:7700/indexes/default/documents' \
  -H 'Content-Type: application/json' -H 'Authorization: Bearer masterKey' \
  --data-binary '[
    {
      "id": 287947,
      "title": "Shazam ⚡️",
      "genres": "comedy"
    }
  ]'
*/

var cli search.Searcher
var index = `default`
var id int64

func TestMain(m *testing.M) {
	defer log.Close()
	id, _ = strconv.ParseInt(time.Now().Format(`20060102150405`), 10, 64)
	var err error
	cli, err = New(search.Config{
		Host:     `http://127.0.0.1:7700`,
		Password: `masterKey`,
	})
	if err != nil {
		panic(err)
	}
	err = cli.InitIndex(&search.IndexConfig{
		Index:      index,
		PrimaryKey: `id`,
	})
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestAdd(t *testing.T) {
	err := cli.Add(index, `id`, map[string]interface{}{
		"title": "钢铁侠美国队长复仇者联盟",
		"id":    id,
	}, map[string]interface{}{
		"title": "Philadelphia",
		"id":    id + 1,
	})
	assert.NoError(t, err)
}

func TestSearch(t *testing.T) {
	count, result, err := cli.Search(index, `钢铁侠`, &search.SearchRequest{
		Limit:  1000,
		Filter: "",
	})
	assert.NoError(t, err)
	assert.Greater(t, count, int64(0))
	com.Dump(result)
	_ = count
}

func TestDelete(t *testing.T) {
	err := cli.Delete(index, fmt.Sprintf(`%d`, id))
	assert.NoError(t, err)
	err = cli.Delete(index, fmt.Sprintf(`%d`, id+1))
	assert.NoError(t, err)
}
