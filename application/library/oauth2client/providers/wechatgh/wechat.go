package wechatgh

import (
	"github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/silenceper/wechat/v2/officialaccount/server"
	"github.com/webx-top/echo"
)

func defaultHandler(svr *server.Server, msg *message.MixMessage) *message.Reply {
	//svr.GetAccessToken()
	// 回复消息：演示回复用户发送的消息
	text := message.NewText(msg.Content)
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
}

func MessageSystem(ctx echo.Context, cfg *config.Config, messageHandler func(*server.Server, *message.MixMessage) *message.Reply) error {
	server := GetServer(ctx, cfg)
	// 设置接收消息的处理方法
	server.SetMessageHandler(func(mm *message.MixMessage) *message.Reply {
		return messageHandler(server, mm)
	})
	// 处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		return err
	}
	// 发送回复的消息
	return server.Send()
}
