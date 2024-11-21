package user

import (
	"context"

	"github.com/admpub/websocket"
	"github.com/coscms/webfront/middleware/sessdata"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

func makeNotice(msgGetter func(context.Context) ([]byte, error)) func(c *websocket.Conn, ctx echo.Context) error {
	return func(c *websocket.Conn, ctx echo.Context) error {
		sessionID := ctx.Session().ID()
		customer := sessdata.Customer(ctx)
		if len(sessionID) > 0 || customer != nil {
			onlineM := modelCustomer.NewOnline(ctx)
			onlineM.SessionId = sessionID
			if customer != nil {
				onlineM.CustomerId = customer.Id
			}
			err := onlineM.Incr(1)
			if err != nil {
				return err
			}
			defer onlineM.Decr(1)
		}
		//push(writer)
		if msgGetter != nil {
			go func() {
				for {
					message, err := msgGetter(ctx)
					if err != nil {
						ctx.Logger().Error(err.Error())
						c.Close()
						return
					}
					c.WriteMessage(websocket.TextMessage, message)
				}
			}()
		}

		//echo
		execute := func(conn *websocket.Conn) error {
			for {
				mt, message, err := conn.ReadMessage()
				if err != nil {
					return err
				}

				// if err = conn.WriteMessage(mt, message); err != nil {
				// 	return err
				// }
				_, _ = mt, message
			}
		}
		err := execute(c)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				ctx.Logger().Debug(err.Error())
			} else {
				ctx.Logger().Error(err.Error())
			}
		}
		return nil
	}
}

func resetClientCount() {
	ctx := defaults.NewMockContext()
	m := modelCustomer.NewOnline(ctx)
	m.ResetClientCount(true)
}
