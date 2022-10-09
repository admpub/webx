package segment_test

import (
	"testing"
	"os"
	"path/filepath"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/testing/test"
	_ "github.com/admpub/webx/application/library/search/segment/gojieba"
	"github.com/admpub/webx/application/library/search/segment"
)

func TestSplitWords(t *testing.T) {
	rootDir := filepath.Join(os.Getenv("GOPATH"),`src`,`github.com/admpub/webx`)
	echo.SetWorkDir(rootDir)
	keywords := `我爱你中国`
	segment.ReloadDict(filepath.Join(rootDir, `data`, `sego`))
	words := segment.SplitWords([]byte(keywords))
	echo.Dump(words)
	expected := []string{`我爱你`,`中国`}
	test.Eq(t,expected,words)
}
