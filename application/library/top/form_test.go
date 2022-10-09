package top

import (
	"testing"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
	"github.com/webx-top/echo/testing/test"
)

func TestParseURLPathKV(t *testing.T) {
	result := ParseURLPathKV(`/a/2/b/3/c/5/d`)
	test.Eq(t, param.StringMap{
		`a`: param.String(`2`),
		`b`: param.String(`3`),
		`c`: param.String(`5`),
		`d`: param.String(``),
	}, result)
}

func TestParseURLPathValues(t *testing.T) {
	result := ParseURLPathValues(`/2_3_5`, `_`, `a`, `b`, `c`, `d`)
	test.Eq(t, param.StringMap{
		`a`: param.String(`2`),
		`b`: param.String(`3`),
		`c`: param.String(`5`),
	}, result)
}

func TestMakeValuesURLPath(t *testing.T) {
	result := MakeValuesURLPath(echo.H{
		`a`: 1,
		`b`: 2,
		`c`: 3,
	}, `_`, `a`, `b`, `c`, `d`)
	test.Eq(t, `1_2_3`, result)
	result = MakeValuesURLPath(echo.H{
		`a`: 1,
		`b`: 2,
		`c`: 3,
		`d`: 4,
	}, `_`, `a`, `b`, `c`, `d`)
	test.Eq(t, `1_2_3_4`, result)
}
