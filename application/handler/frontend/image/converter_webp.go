//go:build webp
// +build webp

package image

import (
	"bytes"
	"io"

	"github.com/admpub/nging/v4/application/registry/upload/convert"
	"github.com/webx-top/image/webp"
)

func init() {
	convert.Register(`.webp`, func(r io.Reader, quality int) (*bytes.Buffer, error) {
		buf, err := webp.Encode(r, float32(quality))
		if err != nil {
			return nil, err
		}
		return buf, err
	})
}
