package aliyun

import (
	"errors"

	app "github.com/KenmyZhang/aliyun-communicate"
	"github.com/coscms/sms"
	"github.com/webx-top/com"
)

var _ sms.Sender = &Aliyun{}

func New() *Aliyun {
	return &Aliyun{
		GatewayURL: `http://dysmsapi.aliyuncs.com/`,
	}
}

type Aliyun struct {
	GatewayURL   string
	AccessKey    string
	AccessSecret string
	SignName     string //默认签名
	TmplCode     string //默认模板代码
	client       *app.SmsClient
}

func (a *Aliyun) Send(c *sms.Config) error {
	tmplCode := c.Template
	signName := c.SignName
	//code := fmt.Sprint(c.Extra(`code`))
	if len(tmplCode) == 0 {
		tmplCode = a.TmplCode
	}
	if len(signName) == 0 {
		signName = a.SignName
	}
	if a.client == nil {
		a.client = app.New(a.GatewayURL)
	}
	b, e := com.JSONEncode(c.ExtraData)
	if e != nil {
		return e
	}
	result, err := a.client.Execute(a.AccessKey, a.AccessSecret, c.Mobile, signName, tmplCode, string(b))
	if err != nil {
		return err
	}
	if result.IsSuccessful() {
		return nil
	}
	if result.Message == `OK` {
		return nil
	}
	return errors.New(`AliyunSMS: ` + result.Message)
}
