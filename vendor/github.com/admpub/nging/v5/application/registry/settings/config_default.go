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
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"

	"github.com/admpub/errors"
	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/dbschema"
)

var configDefaults = map[string]map[string]*dbschema.NgingConfig{
	`base`: {
		`apiKey`: {
			Key:         `apiKey`,
			Label:       `API密钥`,
			Description: ``,
			Value:       ``,
			Group:       `base`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`backendURL`: {
			Key:         `backendURL`,
			Label:       `后台网址`,
			Description: ``,
			Value:       ``,
			Group:       `base`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`debug`: {
			Key:         `debug`,
			Label:       `调试模式`,
			Description: ``,
			Value:       `1`,
			Group:       `base`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`pprof`: {
			Key:         `pprof`,
			Label:       `pprof调优`,
			Description: ``,
			Value:       ``,
			Group:       `base`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`maxRequestBodySize`: {
			Key:         `maxRequestBodySize`,
			Label:       `最大提交尺寸`,
			Description: ``,
			Value:       `2M`,
			Group:       `base`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
	},
	`smtp`: {
		`username`: {
			Key:         `username`,
			Label:       `登录名`,
			Description: ``,
			Value:       ``,
			Group:       `smtp`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`password`: {
			Key:         `password`,
			Label:       `密码`,
			Description: ``,
			Value:       ``,
			Group:       `smtp`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
			Encrypted:   `Y`,
		},
		`host`: {
			Key:         `host`,
			Label:       `服务器`,
			Description: ``,
			Value:       `smtp.exmail.qq.com`,
			Group:       `smtp`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`port`: {
			Key:         `port`,
			Label:       `端口`,
			Description: ``,
			Value:       `465`,
			Group:       `smtp`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`secure`: {
			Key:         `secure`,
			Label:       `认证方式`,
			Description: ``,
			Value:       `SSL`,
			Group:       `smtp`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`identity`: {
			Key:         `identity`,
			Label:       `身份`,
			Description: ``,
			Value:       ``,
			Group:       `smtp`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`timeout`: {
			Key:         `timeout`,
			Label:       `超时时间`,
			Description: ``,
			Value:       ``,
			Group:       `smtp`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`engine`: {
			Key:         `engine`,
			Label:       `发送引擎`,
			Description: ``,
			Value:       `mail`,
			Group:       `smtp`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`from`: {
			Key:         `from`,
			Label:       `发信人地址`,
			Description: ``,
			Value:       ``,
			Group:       `smtp`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`queueSize`: {
			Key:         `queueSize`,
			Label:       `并发数量`,
			Description: ``,
			Value:       `10`,
			Group:       `smtp`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
	},
	`log`: {
		`saveFile`: {
			Key:         `saveFile`,
			Label:       `保存路径`,
			Description: ``,
			Value:       ``,
			Group:       `log`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`fileMaxBytes`: {
			Key:         `fileMaxBytes`,
			Label:       `日志文件尺寸`,
			Description: ``,
			Value:       ``,
			Group:       `log`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
		`targets`: {
			Key:         `targets`,
			Label:       `输出`,
			Description: ``,
			Value:       `console`,
			Group:       `log`,
			Type:        `list`,
			Sort:        0,
			Disabled:    `N`,
		},
		`colorable`: {
			Key:         `colorable`,
			Label:       `彩色日志`,
			Description: ``,
			Value:       `0`,
			Group:       `log`,
			Type:        `text`,
			Sort:        0,
			Disabled:    `N`,
		},
	},
}

func AddDefaultConfig(group string, configs map[string]*dbschema.NgingConfig) {
	if strings.Contains(group, `.`) {
		panic(`Group name is not allowed to contain ".": ` + group)
	}
	if _, y := configDefaults[group]; !y {
		configDefaults[group] = map[string]*dbschema.NgingConfig{}
	}
	for key, conf := range configs {
		if conf.Group != group {
			conf.Group = group
		}
		if conf.Disabled != `N` && conf.Disabled != `Y` {
			conf.Disabled = `N`
		}
		configDefaults[group][key] = conf
	}
}

func DeleteDefaultConfig(group string, keys ...string) {
	if strings.Contains(group, `.`) {
		panic(`Group name is not allowed to contain ".": ` + group)
	}
	if _, y := configDefaults[group]; !y {
		return
	}
	if len(keys) == 0 {
		delete(configDefaults, group)
	} else {
		for _, key := range keys {
			delete(configDefaults[group], key)
		}
		if len(configDefaults[group]) == 0 {
			delete(configDefaults, group)
		}
	}
}

func GetDefaultConfig(group string) map[string]*dbschema.NgingConfig {
	r, _ := configDefaults[group]
	return r
}

func GetDefaultConfigOk(group string) (map[string]*dbschema.NgingConfig, bool) {
	r, y := configDefaults[group]
	return r, y
}

func ConfigHasGroup(group string) bool {
	_, y := configDefaults[group]
	return y
}

func ConfigHasKey(group string, key string) bool {
	g, y := configDefaults[group]
	if !y {
		return false
	}
	_, y = g[key]
	return y
}

func ConfigDefaultsAsStore() echo.H {
	return configAsStore(configDefaults)
}

func ConfigDefaults() map[string]map[string]*dbschema.NgingConfig {
	return configDefaults
}

func Init(ctx echo.Context) error {
	log.Debug(`Initialize the configuration data in the database table`)
	if ctx == nil {
		ctx = defaults.NewMockContext()
	}
	m := dbschema.NewNgingConfig(ctx)
	_, err := m.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`group`).Group(`group`)
	}, 0, -1)
	if err != nil {
		err = errors.WithMessage(err, `Find configuration data`)
		log.Error(err)
	}
	existsList := m.Objects()
	existsIndex := map[string]int{}
	for index, row := range existsList {
		existsIndex[row.Group] = index
	}
	for _, setting := range Settings() {
		group := setting.Group
		if _, ok := existsIndex[group]; ok {
			continue
		}
		gs, ok := GetDefaultConfigOk(group)
		if !ok {
			continue
		}
		for _, conf := range gs {
			conf.SetContext(ctx)
			_, err = conf.EventOFF().Insert()
			if err != nil {
				err = errors.WithMessage(err, `Add configuration data`)
				log.Error(err)
			}
		}
	}
	return err
}

// ConfigAsStore {Group:{Key:ValueObject}}
func ConfigAsStore(ctx echo.Context, groups ...string) echo.H {
	r := echo.H{}
	m := dbschema.NewNgingConfig(ctx)
	cond := db.NewCompounds()
	cond.Add(db.Cond{`disabled`: `N`})
	if len(groups) > 0 {
		if len(groups) > 1 {
			cond.Add(db.Cond{`group`: db.In(groups)})
		} else {
			cond.Add(db.Cond{`group`: groups[0]})
		}
	}
	m.ListByOffset(nil, nil, 0, -1, cond.And())
	rows := m.Objects()
	if len(rows) == 0 {
		return r
	}
	for _, row := range rows {
		row.SetContext(ctx)
		decoder := GetDecoder(row.Group)
		res, err := DecodeConfigValue(row, decoder)
		if err != nil {
			log.Error(`Parsing system settings "`+row.Group+`.`+row.Key, `": `, err)
			continue
		}
		value := res.Get(`ValueObject`, row.Value)
		if _, y := r[row.Group]; !y {
			r[row.Group] = echo.H{row.Key: value}
		} else {
			r.GetStore(row.Group).Set(row.Key, value)
		}
	}
	return r
}

// {Group:{Key:ValueObject}}
func configAsStore(configList map[string]map[string]*dbschema.NgingConfig) echo.H {
	r := echo.H{}
	for group, configs := range configList {
		v := echo.H{}
		decoder := GetDecoder(group)
		for key, row := range configs {
			if row.Disabled == `Y` {
				continue
			}
			res, err := DecodeConfigValue(row, decoder)
			if err != nil {
				log.Error(`Parsing system settings "`+group+`.`+key, `": `, err)
				continue
			}
			value := res.Get(`ValueObject`, row.Value)
			v.Set(key, value)
		}
		r.Set(group, v)
	}
	return r
}
