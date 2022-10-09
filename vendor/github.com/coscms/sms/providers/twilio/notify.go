package twilio

type Notify struct {
	SmsSid        string //"SMXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	SmsStatus     string //"delivered" queued/failed/sent/delivered/undelivered
	MessageStatus string //"delivered"
	To            string //"+15558675310"
	MessageSid    string //"SMXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	AccountSid    string //"AC02b5946ad2c7b9bdced91e8405cc5f3c"
	From          string //"+15017122661",
	ApiVersion    string //"2010-04-01"
}

func (a *Notify) IsSuccess() bool {
	return a.SmsStatus == `delivered` || a.SmsStatus == `sent`
}

func (a *Notify) IsFailure() bool {
	return a.SmsStatus == `undelivered` || a.SmsStatus == `failed`
}
