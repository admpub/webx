package apiutils

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/coscms/webcore/library/restclient"
	"github.com/admpub/resty/v2"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func Submit(ctx echo.Context, apiURL string, formData url.Values, method ...string) (*echo.RawData, error) {
	data := echo.NewData(ctx)
	_, err := SubmitWithRecv(ctx, data, apiURL, formData, method...)
	return data, err
}

func SubmitWithRecv(ctx echo.Context, recv interface{}, apiURL string, formData url.Values, method ...string) (*resty.Response, error) {
	request := restclient.Resty()
	request.SetResult(recv).SetHeader(echo.HeaderAccept, echo.MIMEApplicationJSONCharsetUTF8).SetFormDataFromValues(formData)
	var resp *resty.Response
	var err error
	if len(method) > 0 && strings.EqualFold(method[0], `GET`) {
		resp, err = request.Get(apiURL)
	} else {
		resp, err = request.Post(apiURL)
	}
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf(`%w: %s: %s`, err, apiURL, com.StripTags(resp.String()))
		}
		return nil, fmt.Errorf(`%w: %s`, err, apiURL)
	}
	if !resp.IsSuccess() {
		return resp, fmt.Errorf(`%s: %s: %s`, apiURL, resp.Status(), com.StripTags(resp.String()))
	}
	return resp, err
}

func Submitx(ctx echo.Context, apiURL string, formData url.Values, method ...string) (echo.H, string, error) {
	data := echo.H{}
	var body string
	resp, err := SubmitWithRecv(ctx, &data, apiURL, formData, method...)
	if resp != nil {
		body = com.StripTags(resp.String())
	}
	return data, body, err
}
