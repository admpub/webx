package zincsearch

import (
	"testing"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/webx/application/library/search"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

var cli search.Searcher
var index = `default`
var id string

func TestMain(m *testing.M) {
	defer log.Close()
	id = time.Now().Format(`20060102150405.000`)
	var err error
	cli, err = New(search.Config{
		Host:     `https://playground.dev.zincsearch.com`,
		User:     `admin`,
		Password: `Complexpass#123`,
	})
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestAdd(t *testing.T) {
	err := cli.Add(index, `id`, map[string]string{
		"title": "钢铁侠美国队长复仇者联盟",
		"id":    id,
	})
	assert.NoError(t, err)
}

func TestSearch(t *testing.T) {
	count, result, err := cli.Search(index, `钢铁侠`, &search.SearchRequest{
		SearchType: `matchphrase`,
		Limit:      1000,
	})
	assert.NoError(t, err)
	assert.Greater(t, count, int64(0))
	com.Dump(result)
}

func TestDelete(t *testing.T) {
	err := cli.Delete(index, id)
	assert.NoError(t, err)
}
