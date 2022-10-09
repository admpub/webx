package downloadByContent

import (
	"testing"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/testing/test"
)

func TestOutsideImage(t *testing.T) {
	content := `<img src="http://www.admpub.com/a.jpg" />`
	result := OutsideImage(content, `html`)
	echo.Dump(result)
	test.Eq(t, map[string]string{
		content: `http://www.admpub.com/a.jpg`,
	}, result)
	content = `<img src="http://www.admpub.com/a.jpg" /><img src="http://www.admpub.com/b.png" />`
	result = OutsideImage(content, `html`)
	echo.Dump(result)
	test.Eq(t, map[string]string{
		`<img src="http://www.admpub.com/a.jpg" />`: `http://www.admpub.com/a.jpg`,
		`<img src="http://www.admpub.com/b.png" />`: `http://www.admpub.com/b.png`,
	}, result)

	content = `![test image](http://www.admpub.com/a.jpg)`
	result = OutsideImage(content, `markdown`)
	echo.Dump(result)
	test.Eq(t, map[string]string{
		content: `http://www.admpub.com/a.jpg`,
	}, result)
	content = `![test image](http://www.admpub.com/a.jpg)![test2 image](http://www.admpub.com/b.png "option title")`
	result = OutsideImage(content, `markdown`)
	echo.Dump(result)
	test.Eq(t, map[string]string{
		`![test image](http://www.admpub.com/a.jpg)`:                 `http://www.admpub.com/a.jpg`,
		`![test2 image](http://www.admpub.com/b.png "option title")`: `http://www.admpub.com/b.png`,
	}, result)
	content = `![](https://img.shields.io/github/stars/pandao/editor.md.svg) ![](https://img.shields.io/github/forks/pandao/editor.md.svg) ![](https://img.shields.io/github/tag/pandao/editor.md.svg)`
	result = OutsideImage(content, `markdown`)
	echo.Dump(result)
	test.Eq(t, map[string]string{}, result) // 不支持.svg图片所以捕获为空
	//panic(`ss`)
}
