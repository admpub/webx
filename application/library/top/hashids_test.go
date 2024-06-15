package top

import (
	"fmt"
	"testing"

	"github.com/webx-top/echo/testing/test"
)

func TestGenSN(t *testing.T) {
	sn, err := GenSN(`A`)
	if err != nil {
		panic(err)
	}
	fmt.Println(`sn`, sn) // A3ZJR12RYOO8Q9R
	test.NotEmpty(t, sn)
	test.True(t, IsSN(sn))

	sn, err = snCodec.EncodeUint64([]uint64{12345678901234567890})
	if err != nil {
		panic(err)
	}
	fmt.Println(`sn`, sn)
	test.NotEmpty(t, sn) // 58Z2XMPY803OORY
	test.True(t, IsSN(sn))
}

func TestHideContent1(t *testing.T) {
	rawContent := "[hide]AAA[/hide]efefe[hide]BBB[/hide]22222222\n[hide]CCC[/hide]\n"
	content := HideContent(rawContent, `text`, func(hideType string, hideContent string, args ...string) (hide bool, msgOnHide string) {
		return true, `此处内容需要评论回复后方可阅读`
	}, nil)
	fmt.Println(content)
	test.Eq(t, "[ 此处内容需要评论回复后方可阅读 ]efefe[ 此处内容需要评论回复后方可阅读 ]22222222\n[ 此处内容需要评论回复后方可阅读 ]\n", content)
	content = HideContent(rawContent, `text`, func(hideType string, hideContent string, args ...string) (hide bool, msgOnHide string) {
		return false, ``
	}, nil)
	test.Eq(t, "AAAefefeBBB22222222\nCCC\n", content)
	content = HideContent(`[hide][/hide]`, `text`, func(hideType string, hideContent string, args ...string) (hide bool, msgOnHide string) {
		return false, ``
	}, nil)
	test.Eq(t, "", content)
	content = HideContent(`[hide]A[/hide]`, `text`, func(hideType string, hideContent string, args ...string) (hide bool, msgOnHide string) {
		return false, ``
	}, nil)
	test.Eq(t, "A", content)
	content = HideContent("[hide]\nA\n[/hide]", `text`, func(hideType string, hideContent string, args ...string) (hide bool, msgOnHide string) {
		return false, ``
	}, nil)
	test.Eq(t, "\nA\n", content)
	var hideType2 string
	content = HideContent("[hide:signIn]A[/hide]", `text`, func(hideType string, hideContent string, args ...string) (hide bool, msgOnHide string) {
		hideType2 = hideType
		return false, ``
	}, nil)
	test.Eq(t, "A", content)
	test.Eq(t, `signIn`, hideType2)
}
