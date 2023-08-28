package charset

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/webx-top/com"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var aliases = map[string]string{
	`UTF8`:      `UTF-8`,
	`HZ-GB2312`: `GB2312`,
}

var encodings = map[string]encoding.Encoding{
	`GB18030`: simplifiedchinese.GB18030,
	`GB2312`:  simplifiedchinese.HZGB2312,
	`GBK`:     simplifiedchinese.GBK,
	`UTF-8`:   encoding.Nop,
}

func Supported() []string {
	r := make([]string, 0, len(encodings))
	for k := range encodings {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}

func Register(charset string, encoding encoding.Encoding, alias ...string) {
	charset = strings.ToUpper(charset)
	encodings[charset] = encoding
	for _, a := range alias {
		a = strings.ToUpper(a)
		aliases[a] = charset
	}
}

func Encoding(charset string) encoding.Encoding {
	charset = strings.ToUpper(charset)
	if cs, ok := aliases[charset]; ok {
		charset = cs
	}
	if enc, ok := encodings[charset]; ok {
		return enc
	}
	return nil
}

func NewTransformWriter(charset string, dst io.WriteCloser) (io.WriteCloser, error) {
	cs := Encoding(charset)
	if nil == cs {
		return nil, fmt.Errorf("charset '%s' is unsupported", charset)
	}
	if cs == encoding.Nop {
		return dst, nil
	}
	return transform.NewWriter(dst, cs.NewDecoder()), nil
}

func NewTransformReader(charset string, src io.Reader) (io.Reader, error) {
	cs := Encoding(charset)
	if nil == cs {
		return nil, fmt.Errorf("charset '%s' is unsupported", charset)
	}
	if cs == encoding.Nop {
		return src, nil
	}
	return transform.NewReader(src, cs.NewDecoder()), nil
}

func Transform(charset string, content string) (string, error) {
	r := strings.NewReader(content)
	tr, err := NewTransformReader(charset, r)
	if err != nil {
		return content, err
	}
	b, err := io.ReadAll(tr)
	if err != nil {
		return content, err
	}
	return com.Bytes2str(b), nil
}

func TransformBytes(charset string, content []byte) ([]byte, error) {
	r := bytes.NewReader(content)
	tr, err := NewTransformReader(charset, r)
	if err != nil {
		return content, err
	}
	b, err := io.ReadAll(tr)
	if err != nil {
		return content, err
	}
	return b, nil
}

func NewTransformFunc(charset string) (func(string) (string, error), error) {
	cs := Encoding(charset)
	if nil == cs {
		return nil, fmt.Errorf("charset '%s' is unsupported", charset)
	}
	if cs == encoding.Nop {
		return func(v string) (string, error) { return v, nil }, nil
	}
	t := cs.NewDecoder()
	return func(content string) (string, error) {
		r := strings.NewReader(content)
		tr := transform.NewReader(r, t)
		b, err := io.ReadAll(tr)
		if err != nil {
			return content, err
		}
		return com.Bytes2str(b), nil
	}, nil
}

func NewTransformBytesFunc(charset string) (func([]byte) ([]byte, error), error) {
	cs := Encoding(charset)
	if nil == cs {
		return nil, fmt.Errorf("charset '%s' is unsupported", charset)
	}
	if cs == encoding.Nop {
		return func(v []byte) ([]byte, error) { return v, nil }, nil
	}
	t := cs.NewDecoder()
	return func(content []byte) ([]byte, error) {
		r := bytes.NewReader(content)
		tr := transform.NewReader(r, t)
		b, err := io.ReadAll(tr)
		if err != nil {
			return content, err
		}
		return b, nil
	}, nil
}
