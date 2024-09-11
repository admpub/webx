package xsocketio

import (
	"net/http"
	"strings"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/common"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/webx-top/echo"
	esi "github.com/webx-top/echo-socket.io"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware"
)

func RegisterRoute(e echo.RouteRegister, s ...func(*middleware.CORSConfig)) {
	cfg := &middleware.CORSConfig{
		AllowOrigins: []string{`Token`},
	}
	for _, f := range s {
		f(cfg)
	}
	prefix := e.Prefix()
	nsp := strings.Trim(prefix, `/`)
	nsp = strings.ReplaceAll(nsp, `/`, `_`)
	socket := SocketIO(nsp)
	e.Any(`/socket.io/`, func(ctx echo.Context) error {
		if common.Setting(`socketio`).String(`enabled`) != `1` {
			return echo.ErrNotFound
		}
		return socket.Handle(ctx)
	}, middleware.CORSWithConfig(*cfg))
}

var (
	events       = []func(esi.IWrapper){}
	onConnect    = []func(ctx echo.Context, conn socketio.Conn) error{}
	onError      = []func(ctx echo.Context, conn socketio.Conn, e error){}
	onDisconnect = []func(ctx echo.Context, conn socketio.Conn, msg string){}

	RequestChecker engineio.CheckerFunc = func(req *http.Request) (http.Header, error) {
		token := common.Setting(`socketio`).String(`token`)
		if len(token) == 0 {
			return nil, nil
		}
		post := req.Header.Get(`Token`)
		if len(post) == 0 {
			post = req.URL.Query().Get(`token`)
		}
		if token != post {
			if log.IsEnabled(log.LevelDebug) {
				log.Debugf(`[socketIO] invalid token: %q`, post)
				log.Debugf(`[socketIO] request headers: %+v`, req.Header)
			}
			return nil, echo.NewError(`invalid token`, code.InvalidToken)
		}
		return nil, nil
	}
)

func OnEvent(fns ...func(esi.IWrapper)) {
	events = append(events, fns...)
}

func OnConnect(fns ...func(ctx echo.Context, conn socketio.Conn) error) {
	onConnect = append(onConnect, fns...)
}

func OnError(fns ...func(ctx echo.Context, conn socketio.Conn, e error)) {
	onError = append(onError, fns...)
}

func OnDisconnect(fns ...func(ctx echo.Context, conn socketio.Conn, msg string)) {
	onDisconnect = append(onDisconnect, fns...)
}

func socketIOWrapper(nsp string) *esi.Wrapper {
	wrapper := esi.NewWrapper(&engineio.Options{
		RequestChecker: RequestChecker,
	})

	wrapper.OnConnect(nsp, func(ctx echo.Context, conn socketio.Conn) error {
		for _, fn := range onConnect {
			if err := fn(ctx, conn); err != nil {
				return err
			}
		}
		return nil
	})

	wrapper.OnError(nsp, func(ctx echo.Context, conn socketio.Conn, e error) {
		log.Error("[socketIO] meet error: ", e)
		for _, fn := range onError {
			fn(ctx, conn, e)
		}
		conn.Close()
	})

	wrapper.OnDisconnect(nsp, func(ctx echo.Context, conn socketio.Conn, msg string) {
		log.Debug("[socketIO] closed", msg)
		for _, fn := range onDisconnect {
			fn(ctx, conn, msg)
		}
		conn.Close()
	})

	for _, fn := range events {
		fn(wrapper)
	}

	return wrapper
}
