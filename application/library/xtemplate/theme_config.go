package xtemplate

import (
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/admpub/confl"
	"github.com/coscms/forms"
	formCommon "github.com/coscms/forms/common"
	formConfig "github.com/coscms/forms/config"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	_ "github.com/admpub/nging/v5/application/library/formbuilder"
	_ "github.com/admpub/webx/application/library/formbuilder"
	//_ "github.com/coscms/forms/defaults"
)

var (
	themeNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

func IsThemeName(s string) bool {
	return themeNameRegex.MatchString(s)
}

func NewThemeInfo() *ThemeInfo {
	return &ThemeInfo{
		CustomConfig: echo.H{},
		FormConfig:   make(map[string]formConfig.Config),
	}
}

type ThemeAuthor struct {
	Name        string `json:"name,omitempty"`        // 作者网名
	Email       string `json:"email,omitempty"`       // 作者E-mail地址
	URL         string `json:"url,omitempty"`         // 作者网址
	Description string `json:"description,omitempty"` // 作者简介
}

type ThemeColors []ThemeColor

func (t ThemeColors) HasName(colorName string) bool {
	for _, ti := range t {
		if ti.Name == colorName {
			return true
		}
	}
	return false
}

type ThemeColor struct {
	Name         string `json:"name,omitempty"`         // 颜色英文名(用于调用相应颜色的css文件)
	Title        string `json:"title,omitempty"`        // 颜色中文标题
	IsDefault    bool   `json:"isDefault,omitempty"`    // 是否为默认颜色
	Color        string `json:"color,omitempty"`        // 颜色值
	PreviewImage string `json:"previewImage,omitempty"` // 预览图
}

func RGB2Hex(rgb string) string {
	if len(rgb) == 0 {
		return ``
	}
	if strings.HasPrefix(rgb, `#`) {
		return rgb
	}
	n := 3
	if strings.HasPrefix(rgb, `rgba(`) {
		rgb = strings.TrimPrefix(rgb, `rgba(`)
		n = 4
	} else {
		rgb = strings.TrimPrefix(rgb, `rgb(`)
	}
	rgb = strings.TrimRight(rgb, `);`)
	parts := strings.SplitN(rgb, `,`, n)
	if len(parts) < 3 {
		return ``
	}
	rgbNums := make([]int, 3)
	for x := 0; x < 3; x++ {
		s := strings.TrimSpace(parts[x])
		var i int
		if len(s) > 0 {
			i, _ = strconv.Atoi(s)
			if i > 255 || i < 0 {
				return ``
			}
		}
		rgbNums[x] = i
	}
	return fmt.Sprintf(`#%X%X%X`, rgbNums[0], rgbNums[1], rgbNums[2])
}

func (t ThemeColor) HexColor() string {
	return RGB2Hex(t.Color)
}

// ThemeInfo 模板主题信息
type ThemeInfo struct {
	Author       ThemeAuthor                  `json:"author,omitempty"`       // 模板作者
	Colors       ThemeColors                  `json:"colors,omitempty"`       // 多种颜色主题时，所有支持的颜色信息
	Version      string                       `json:"version,omitempty"`      // 版本号(格式: 1.0.0)
	Name         string                       `json:"name"`                   // 主题英文名
	Title        string                       `json:"title"`                  // 主题中文标题
	UpdatedAt    string                       `json:"updatedAt,omitempty"`    // 更新时间(格式: 2006-01-02 15:04:05)
	Description  string                       `json:"description,omitempty"`  // 简介
	PreviewImage string                       `json:"previewImage,omitempty"` // 预览图
	PreviewURL   string                       `json:"previewURL,omitempty"`   // 预览网址
	CustomConfig echo.H                       `json:"customConfig,omitempty"` // 可自定义配置的数据
	FormConfig   map[string]formConfig.Config `json:"formConfig,omitempty"`   // 表单配置
	Fallback     []string                     `json:"fallback,omitempty"`     // 兜底主题
	embed        bool
}

type ThemeInfoLite struct {
	Name         string `json:"name"`         // 主题英文名
	Title        string `json:"title"`        // 主题中文标题
	PreviewImage string `json:"previewImage"` // 预览图
}

func (t *ThemeInfo) AsLite() ThemeInfoLite {
	return ThemeInfoLite{
		Name:         t.Name,
		Title:        t.Title,
		PreviewImage: t.PreviewImage,
	}
}

func (t *ThemeInfo) HasColorName(colorName string) bool {
	return t.Colors.HasName(colorName)
}

func (t *ThemeInfo) HasForm(templateName string) bool {
	if len(t.Colors) > 0 {
		return true
	}
	_, ok := t.FormConfig[templateName]
	return ok
}

func (t *ThemeInfo) SetEmbed() *ThemeInfo {
	t.embed = true
	return t
}

func (t *ThemeInfo) Embed() bool {
	return t.embed
}

func (t *ThemeInfo) SaveForm(ctx echo.Context, templateName string, gets ...func(fieldName string, fieldValue string) error) error {
	cfg, ok := t.FormConfig[templateName]
	if !ok {
		return nil
	}
	var get func(fieldName string, fieldValue string) error
	if len(gets) > 0 {
		get = gets[0]
	}
	if get == nil {
		get = func(fieldName string, fieldValue string) error {
			t.CustomConfig.Set(fieldName, ctx.Form(fieldName))
			return nil
		}
	}
	return cfg.GetValue(get)
}

func (t *ThemeInfo) Render(templateName string) template.HTML {
	cfg, ok := t.FormConfig[templateName]
	if !ok {
		return param.EmptyHTML
	}
	copied := cfg.Clone()
	copied.SetValue(func(fieldName string) string {
		return t.CustomConfig.String(fieldName)
	})
	f := forms.NewForms(forms.New())
	if len(cfg.Theme) == 0 {
		copied.Theme = formCommon.BOOTSTRAP
	}
	f.Theme = formCommon.BOOTSTRAP
	f.Init(copied)
	f.ParseFromConfig(true)
	f.Config().Template = formCommon.TmplDir(f.Config().Theme) + `/allfields.html`
	return f.Render()
}

func (t *ThemeInfo) Encode() []byte {
	b, _ := confl.Marshal(t)
	return b
}

func (t *ThemeInfo) EncodeToFile(destFile string) error {
	b, err := confl.Marshal(t)
	if err != nil {
		return err
	}
	err = os.WriteFile(destFile, b, os.ModePerm)
	return err
}

func (t *ThemeInfo) fixedColors() {
	for i, v := range t.Colors {
		v.Name = strings.TrimSpace(v.Name)
		v.Color = strings.TrimSpace(v.Color)
		t.Colors[i] = v
	}
}

func (t *ThemeInfo) Decode(b []byte) error {
	err := confl.Unmarshal(b, t)
	if err == nil {
		t.fixedColors()
	}
	return err
}

func (t *ThemeInfo) DecodeFile(file string) error {
	_, err := confl.DecodeFile(file, t)
	if err == nil {
		t.fixedColors()
	}
	return err
}

func GetThemeInfoFromContext(ctx echo.Context) *ThemeInfo {
	v, _ := ctx.Internal().Get(`theme.config`).(*ThemeInfo)
	return v
}
