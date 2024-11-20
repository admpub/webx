package user

import (
	"github.com/admpub/websocket"
	"github.com/coscms/webfront/middleware/sessdata"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/echo"
)

func Notice(c *websocket.Conn, ctx echo.Context) error {
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
	// go func() {
	// 	for {
	// 		c.WriteMessage(websocket.TextMessage, msgBytes)
	// 	}
	// }()

	//echo
	var execute = func(conn *websocket.Conn) error {
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				return err
			}

			if err = conn.WriteMessage(mt, message); err != nil {
				return err
			}
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
