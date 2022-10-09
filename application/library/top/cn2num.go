package top

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 一百二十三 => 123
// 一百 => 100
// 二十 => 20
// 十三 => 13
var nRegex = regexp.MustCompile(`[零一二三四五六七八九]+`)
var pRegex = regexp.MustCompile(`[十百千万亿]+`)
var splitRegex = regexp.MustCompile(`[万亿]`)
var numDict = map[string]uint{
	`零`: 0,
	`一`: 1,
	`二`: 2,
	`三`: 3,
	`四`: 4,
	`五`: 5,
	`六`: 6,
	`七`: 7,
	`八`: 8,
	`九`: 9,
}

var postDict = map[string]uint64{
	`十`: 10,
	`百`: 100,
	`千`: 1000,
	`万`: 10000,
	`亿`: 100000000,
}

func CN2Number(cn string) (uint64, error) {
	if cn == `十` {
		return 10, nil
	}
	cn = nRegex.ReplaceAllStringFunc(cn, func(repl string) string {
		return fmt.Sprint(numDict[repl])
	})
	if strings.HasPrefix(cn, `十`) {
		cn = strings.TrimPrefix(cn, `十`)
		i, err := strconv.ParseInt(`1`+cn, 10, 64)
		return uint64(i), err
	}
	for key, value := range postDict {
		if strings.HasSuffix(cn, key) {
			cn = pRegex.ReplaceAllString(cn, ``)
			i, err := strconv.ParseInt(cn, 10, 64)
			if err != nil {
				return uint64(i), err
			}
			return uint64(i) * value, nil
		}
	}
	cn = pRegex.ReplaceAllString(cn, ``)
	i, err := strconv.ParseInt(cn, 10, 64)
	return uint64(i), err
}
