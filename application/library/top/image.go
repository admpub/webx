package top

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/subdomains"

	uploadChecker "github.com/coscms/webcore/registry/upload/checker"
)

// TokenURL 带签名网址
func TokenURL(values ...interface{}) string {
	var urlValues url.Values
	if len(values) == 1 {
		switch t := values[0].(type) {
		case url.Values:
			urlValues = t
		case map[string][]string:
			urlValues = url.Values(t)
		default:
			urlValues = tplfunc.URLValues(values...)
		}
	} else {
		urlValues = tplfunc.URLValues(values...)
	}
	unixtime := fmt.Sprint(time.Now().Unix())
	urlValues.Set(`time`, unixtime)
	urlValues.Del(`token`)
	return `/` + uploadChecker.Token(values...) + `?` + urlValues.Encode()
}

func PictureWithDefaultHTML(picURL string, defaultURL string, widthAndHeights ...string) template.HTML {
	if len(picURL) == 0 {
		picURL = defaultURL
		s := `<picture>`
		s += `<img src="` + picURL + `"/>`
		s += `</picture>`
		return template.HTML(s)
	}
	return PictureHTML(picURL, widthAndHeights...)
}

func PictureHTML(picURL string, widthAndHeights ...string) template.HTML {
	s := `<picture>`
	s += `<source type="image/webp" srcset="` + picURL + `.webp">`
	for _, wh := range widthAndHeights {
		s += `<source srcset="` + ImageProxyURL(`type`, ``, `size`, wh, `image`, picURL) + `" media="(max-width: ` + strings.SplitN(wh, `x`, 2)[0] + `px)">`
	}
	s += `<img src="` + picURL + `"/>`
	s += `</picture>`
	return template.HTML(s)
}

func URLFor(purl string) string {
	return subdomains.Default.URL(purl, `frontend`)
}

func AbsoluteURL(purl string) string {
	if !com.IsFullURL(purl) {
		return URLFor(purl)
	}
	return purl
}

// ResizeImageURL ResizeImageURL(imageURL,`1000x1240`,`default.jpg`)
func ResizeImageURL(image string, size string, defaults ...string) string {
	if len(image) == 0 {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return image
	}
	return ImageProxyURL(url.Values{
		`size`:  {size},
		`image`: {image},
	})
}

// ImageProxyURL 图片代理网址
func ImageProxyURL(values ...interface{}) string {
	var urlValues url.Values
	if len(values) == 1 {
		switch t := values[0].(type) {
		case url.Values:
			urlValues = t
		case map[string][]string:
			urlValues = url.Values(t)
		default:
			urlValues = tplfunc.URLValues(values...)
		}
	} else {
		urlValues = tplfunc.URLValues(values...)
	}
	//unixtime := fmt.Sprint(time.Now().Unix())
	//urlValues.Set(`time`, unixtime)
	urlValues.Del(`token`)
	return URLFor(`/image/proxy`) + `/` + uploadChecker.Token(urlValues) + `?` + urlValues.Encode()
}
