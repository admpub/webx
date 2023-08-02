package alert

import (
	"encoding/json"
	"strings"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/cron/send"
	"github.com/admpub/nging/v5/application/library/imbot"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

// Send 发送警报
func (alertData *AlertData) Send(a *dbschema.NgingAlertRecipient) (err error) {
	title := alertData.Title
	ct := alertData.Content
	if ct == nil {
		return
	}
	params := alertData.Data
	switch a.Type {
	case `email`:
		content := ct.EmailContent(params)
		if len(content) == 0 {
			return
		}
		err = send.Mail(a.Account, strings.SplitN(a.Account, `@`, 2)[0], title, content)
	case `webhook`:
		if a.Platform == RecipientPlatformWebhookCustom { // 自定义webhook
			custom := &WebhookCustom{}
			extraBytes := com.Str2bytes(a.Extra)
			if err := json.Unmarshal(extraBytes, custom); err != nil {
				err = common.JSONBytesParseError(err, extraBytes)
				return err
			}
			if len(custom.Url) == 0 {
				if len(a.Account) > 7 {
					switch a.Account[0:7] {
					case `https:/`, `http://`:
						custom.Url = a.Account
					}
				}
			}
			var tagReplacer func(string) string
			var tagValues map[string]string
			if tv, ok := ct.(TagValues); ok {
				tagValues = tv.TagValues()
			} else {
				tagValues = map[string]string{
					`title`:           alertData.Title,
					`emailContent`:    string(ct.EmailContent(params)),
					`markdownContent`: string(ct.MarkdownContent(params)),
				}
			}
			if tagValues != nil {
				tagReplacer = func(content string) string {
					for tag, value := range tagValues {
						content = strings.ReplaceAll(content, `{{`+tag+`}}`, value)
					}
					return content
				}
			}
			return custom.ToWebhook().Exec(tagReplacer, tagReplacer)
		}

		// 内置 imbot 的 webhook
		mess := imbot.Open(a.Platform)
		if mess == nil || mess.Messager == nil {
			return
		}

		content := ct.MarkdownContent(params)
		if len(content) == 0 {
			return
		}
		message := string(content)
		var apiURL string
		if len(a.Account) > 7 {
			switch a.Account[0:7] {
			case `https:/`, `http://`:
				apiURL = a.Account
			}
		}
		if len(apiURL) == 0 {
			apiURL = mess.Messager.BuildURL(a.Account)
		}
		var atMobiles []string
		extra := echo.H{}
		if len(a.Extra) > 0 {
			json.Unmarshal(com.Str2bytes(a.Extra), &extra)
		}
		if extra.Has(`at`) {
			switch v := extra.Get(`at`).(type) {
			case []interface{}:
				atMobiles = make([]string, len(v))
				for k, m := range v {
					atMobiles[k] = param.AsString(m)
				}
			case []string:
				atMobiles = v
			}
		}
		go func(apiURL string, title string, message string, atMobiles ...string) {
			err = mess.Messager.SendMarkdown(apiURL, title, message, atMobiles...)
		}(apiURL, title, message, atMobiles...)
	}
	return
}

func (alertData *AlertData) SendTopic(ctx echo.Context, topic string) (err error) {
	rows := []*AlertTopicExt{}
	m := dbschema.NewNgingAlertTopic(ctx)
	_, err = m.ListByOffset(&rows, nil, 0, -1, db.And(
		db.Cond{`topic`: topic},
		db.Cond{`disabled`: `N`},
	))
	if err != nil {
		return
	}
	for _, row := range rows {
		if row.Recipient != nil && row.Recipient.Disabled == `N` {
			err = alertData.Send(row.Recipient)
		}
	}
	return
}
