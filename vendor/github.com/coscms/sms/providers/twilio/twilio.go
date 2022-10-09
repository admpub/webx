package twilio

import (
	"errors"
	"fmt"

	"github.com/admpub/gotwilio"

	"github.com/coscms/sms"
	"github.com/webx-top/com"
)

var _ sms.Sender = &Twilio{}

func New() *Twilio {
	return &Twilio{}
}

type Twilio struct {
	AccessKey    string
	AccessSecret string
	From         string
	CountryCode  string //+86
	client       *gotwilio.Twilio
}

func (a *Twilio) Send(c *sms.Config) error {
	if a.client == nil {
		a.client = gotwilio.NewTwilioClient(a.AccessKey, a.AccessSecret)
	}
	to := c.Mobile
	if len(a.CountryCode) == 0 {
		to = `+86` + to
	}
	resp, exception, err := a.client.SendSMS(a.From, to, c.Content, "", "")
	if err != nil {
		return err
	}
	if exception != nil {
		err = errors.New(exception.Message + ` (` + fmt.Sprint(exception.Code) + `)`)
	}
	com.Dump(map[string]interface{}{`resp`: resp, `exception`: exception, `err`: err})
	return err
}
