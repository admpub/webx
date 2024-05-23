package gh

import (
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/silenceper/wechat/v2/officialaccount/server"
	"github.com/webx-top/echo"
)

type EventType string

const (
	EventUserGetCard     EventType = `user_get_card`
	EventUserDelCard     EventType = `user_del_card`
	EventUserConsumeCard EventType = `user_consume_card`
	EventSubscribe       EventType = `subscribe`
	EventUnsubscribe     EventType = `unsubscribe`
	EventScan            EventType = `SCAN`
	EventDefault         EventType = ``
)

type Handler func(c echo.Context, s *server.Server, m *message.MixMessage) *message.Reply

var eventHandlers = map[EventType]Handler{}

func RegisterEventHandler(event EventType, handler Handler) {
	eventHandlers[event] = handler
}

func MakeMessageHandler(c echo.Context, msgHandler func(s *server.Server, m *message.MixMessage) *message.Reply) func(*server.Server, *message.MixMessage) *message.Reply {
	return func(s *server.Server, m *message.MixMessage) *message.Reply {
		if m.MsgType != message.MsgTypeEvent {
			if msgHandler == nil {
				return nil
			}
			return msgHandler(s, m)
		}
		event := EventType(m.Event)
		if h, ok := eventHandlers[event]; ok {
			return h(c, s, m)
		}
		if h, ok := eventHandlers[EventDefault]; ok {
			return h(c, s, m)
		}
		return nil
		// 回复消息：演示回复用户发送的消息
		//text := message.NewText(m.Content)
		//return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	}
}
