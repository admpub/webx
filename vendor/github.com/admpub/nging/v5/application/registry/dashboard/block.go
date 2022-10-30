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

package dashboard

import (
	"database/sql"

	"github.com/webx-top/echo"
)

func NewBlock(content func(echo.Context) error) *Block {
	return &Block{content: content}
}

type IsHidden interface {
	IsHidden(ctx echo.Context) (hidden bool)
}

type Block struct {
	Title   string `json:",omitempty" xml:",omitempty"` // 标题
	Ident   string `json:",omitempty" xml:",omitempty"` // 英文标识
	Extra   echo.H `json:",omitempty" xml:",omitempty"` // 附加数据
	Hidden  sql.NullBool
	Tmpl    string //模板文件
	Footer  string //页脚模板文件
	content func(echo.Context) error
	hidden  func(echo.Context) bool
}

func (c *Block) Ready(ctx echo.Context) error {
	if c.content != nil {
		return c.content(ctx)
	}
	return nil
}

func (c *Block) SetTitle(title string) *Block {
	c.Title = title
	return c
}

func (c *Block) SetIdent(ident string) *Block {
	c.Ident = ident
	return c
}

func (c *Block) SetExtra(extra echo.H) *Block {
	c.Extra = extra
	return c
}

func (c *Block) IsHidden(ctx echo.Context) (hidden bool) {
	v, ok := ctx.Internal().GetOk(`registry.dashboard.block.` + c.Ident)
	if ok {
		hidden = v.(bool)
		return
	}
	if c.hidden != nil {
		hidden = c.hidden(ctx)
	} else {
		hidden = c.Hidden.Bool
	}
	ctx.Internal().Set(`registry.dashboard.block.`+c.Ident, hidden)
	return
}

func (c *Block) SetHidden(hidden bool) *Block {
	c.Hidden.Bool = hidden
	c.Hidden.Valid = true
	return c
}

func (c *Block) SetExtraKV(key string, value interface{}) *Block {
	if c.Extra == nil {
		c.Extra = echo.H{}
	}
	c.Extra.Set(key, value)
	return c
}

func (c *Block) SetTmpl(tmpl string) *Block {
	c.Tmpl = tmpl
	return c
}

func (c *Block) SetFooter(footer string) *Block {
	c.Footer = footer
	return c
}

func (c *Block) SetContentGenerator(content func(echo.Context) error) *Block {
	c.content = content
	return c
}

func (c *Block) SetHiddenDetector(hidden func(echo.Context) bool) *Block {
	c.hidden = hidden
	return c
}

type Blocks []*Block

func (c *Blocks) Ready(ctx echo.Context) error {
	for _, blk := range *c {
		if blk != nil {
			if err := blk.Ready(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

// Remove 删除元素
func (c *Blocks) Remove(index int) {
	if index < 0 {
		*c = (*c)[0:0]
		return
	}
	size := c.Size()
	if size > index {
		if size > index+1 {
			*c = append((*c)[0:index], (*c)[index+1:]...)
		} else {
			*c = (*c)[0:index]
		}
	}
}

// Add 添加列表项
func (c *Blocks) Add(index int, list ...*Block) {
	if len(list) == 0 {
		return
	}
	if index < 0 {
		*c = append(*c, list...)
		return
	}
	size := c.Size()
	if size > index {
		list = append(list, (*c)[index])
		(*c)[index] = list[0]
		if len(list) > 1 {
			c.Add(index+1, list[1:]...)
		}
		return
	}
	for start, end := size, index-1; start < end; start++ {
		*c = append(*c, nil)
	}
	*c = append(*c, list...)
}

// Set 设置元素
func (c *Blocks) Set(index int, list ...*Block) {
	if len(list) == 0 {
		return
	}
	if index < 0 {
		*c = append(*c, list...)
		return
	}
	size := c.Size()
	if size > index {
		(*c)[index] = list[0]
		if len(list) > 1 {
			c.Set(index+1, list[1:]...)
		}
		return
	}
	for start, end := size, index-1; start < end; start++ {
		*c = append(*c, nil)
	}
	*c = append(*c, list...)
}

func (c *Blocks) Size() int {
	return len(*c)
}

func (c *Blocks) Search(cb func(*Block) bool) int {
	for index, block := range *c {
		if cb(block) {
			return index
		}
	}
	return -1
}

func (c *Blocks) FindTmpl(tmpl string) int {
	return c.Search(func(block *Block) bool {
		return block.Tmpl == tmpl
	})
}

func (c *Blocks) RemoveByTmpl(tmpl string) {
	index := c.FindTmpl(tmpl)
	if index > -1 {
		c.Remove(index)
	}
}

var blocks = Blocks{}

func BlockRegister(block ...*Block) {
	blocks.Add(-1, block...)
}

func BlockAdd(index int, block ...*Block) {
	blocks.Add(index, block...)
}

// BlockRemove 删除元素
func BlockRemove(index int) {
	blocks.Remove(index)
}

// BlockSet 设置元素
func BlockSet(index int, list ...*Block) {
	blocks.Set(index, list...)
}

func BlockAll(ctx echo.Context) Blocks {
	result := make(Blocks, len(blocks))
	for k, v := range blocks {
		val := *v
		val.SetHidden(val.IsHidden(ctx))
		result[k] = &val
	}
	return result
}
