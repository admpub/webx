package top

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/coscms/webcore/library/common"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func CopyFormDataBy(ctx echo.Context, formData map[string][]string) {
	for key, vals := range formData {
		for idx, val := range vals {
			if idx == 0 {
				ctx.Request().Form().Set(key, val)
			} else {
				ctx.Request().Form().Add(key, val)
			}
		}
	}
}

// ParseURLPathKV 解析网址路径表示的键值数据 /键1/值1/键2/值2
func ParseURLPathKV(path string) param.StringMap {
	path = strings.TrimPrefix(path, `/`)
	args := strings.Split(path, `/`)
	result := param.StringMap{}
	var key string
	for i, j := 0, len(args); i < j; i++ {
		if i%2 == 0 {
			key = args[i]
		} else {
			result[key] = param.String(args[i])
			key = ``
		}
	}
	if len(key) > 0 {
		result[key] = param.String(``)
	}
	return result
}

// ParseURLPathValues 解析网址路径表示的值数据 /值1_值2
// 通过传入参数keys来指定各个对应位置的key
func ParseURLPathValues(path string, sep string, keys ...string) param.StringMap {
	path = strings.TrimPrefix(path, `/`)
	args := strings.Split(path, sep)
	size := len(args)
	result := param.StringMap{}
	for index, key := range keys {
		if index >= size {
			break
		}
		result[key] = param.String(args[index])
	}
	return result
}

// MakeValuesURLPath 生成带参数值的网址路径
func MakeValuesURLPath(args echo.H, sep string, keys ...string) string {
	values := make([]string, len(keys))
	for index, key := range keys {
		values[index] = url.QueryEscape(args.String(key))
	}
	return strings.TrimRight(strings.Join(values, sep), sep)
}

func MakeListURL(urlPrefix string, params param.StringMap,
	filterParamNames []string, sep string, args ...interface{}) string {
	values := common.HPoolGet()
	defer common.HPoolRelease(values)
	for _, key := range filterParamNames {
		values[key] = params.String(key)
	}
	var k string
	for i, j := 0, len(args); i < j; i++ {
		if i%2 == 0 {
			k = fmt.Sprint(args[i])
			continue
		}
		values[k] = args[i]
		k = ``
	}
	if len(k) > 0 {
		values[k] = ``
		k = ``
	}
	return urlPrefix + MakeValuesURLPath(values, sep, filterParamNames...)
}
