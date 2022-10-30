/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package settings

import (
	"errors"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
)

type Encoder func(v *dbschema.NgingConfig, formDataMap echo.H) ([]byte, error)

var encoders = map[string]Encoder{}

func Encoders() map[string]Encoder {
	return encoders
}

func GetEncoder(group string) Encoder {
	ps, _ := encoders[group]
	return ps
}

// RegisterEncoder 注册配置值编码器（用于客户端提交表单数据之后的编码操作，编码结果保存到数据库）
// 名称支持"group"或"group.key"两种格式，例如:
// settings.RegisterDecoder(`sms`,...)对整个sms组的配置有效
// settings.RegisterDecoder(`sms.twilio`,...)对sms组内key为twilio的配置有效
func RegisterEncoder(group string, encoder Encoder) {
	encoders[group] = encoder
}

var ErrNotExists = errors.New(`Not exists`)

func EncodeConfigValue(_v *echo.Mapx, v *dbschema.NgingConfig, encoder Encoder) (value string, err error) {
	if _v.IsMap() {
		var b []byte
		store := _v.AsStore()
		if subEncoder := GetEncoder(v.Group + `.` + v.Key); subEncoder != nil {
			b, err = subEncoder(v, store)
		} else if encoder != nil {
			b, err = encoder(v, store)
		} else {
			b, err = com.JSONEncode(store)
		}
		if err != nil {
			return
		}
		value = string(b)
	} else if _v.IsSlice() {
		items := []string{}
		for _, item := range _v.AsFlatSlice() {
			item = strings.TrimSpace(item)
			if len(item) == 0 {
				return
			}
			items = append(items, item)
		}
		if subEncoder := GetEncoder(v.Group + `.` + v.Key); subEncoder != nil {
			var b []byte
			b, err = subEncoder(v, echo.H{`value`: items})
			if err != nil {
				return
			}
			value = string(b)
		} else {
			value = strings.Join(items, `,`)
		}
	} else {
		if subEncoder := GetEncoder(v.Group + `.` + v.Key); subEncoder != nil {
			var b []byte
			b, err = subEncoder(v, echo.H{`value`: _v.Value()})
			if err != nil {
				return
			}
			value = string(b)
		} else {
			value = _v.Value() //c.Form(group + `[` + v.Key + `]`)
		}
	}
	value = DefaultEncoder(v, value)
	return
}

func DefaultEncoder(v *dbschema.NgingConfig, value string) string {
	switch v.Type {
	case `html`, `text`:
		// 配置数据为后台输入，不用过滤XSS
	default:
		value = common.ContentEncode(value, v.Type)
	}
	if v.Encrypted == `Y` {
		value = echo.Get(common.ConfigName).(Codec).Encode(value)
	}
	return value
}
