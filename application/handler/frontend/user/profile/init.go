package profile

import (
	"github.com/webx-top/echo"

	_ "github.com/admpub/webx/application/handler/frontend/user/wallet"
	"github.com/coscms/webfront/initialize/frontend"
)

func init() {
	frontend.RegisterToGroup(`/user`, func(u echo.RouteRegister) {
		p := u.Group(`/profile`)
		p.Route(`GET,POST`, `/detail`, Profile)
		p.Route(`GET,POST`, `/settings`, Settings)
		p.Route(`GET,POST`, `/password`, Password)
		p.Route(`GET,POST`, `/partial_authentication`, _Authentication)
		p.Route(`GET,POST`, `/binding`, Binding)
		p.Route(`POST`, `/follow`, Follow)
		p.Route(`POST`, `/unfollow`, Unfollow)
		p.Route(`GET,POST`, `/is_followed`, IsFollowed)
		p.Route(`GET,POST`, `/following`, Following)
		p.Route(`GET,POST`, `/followers`, Followers)
		p.Route(`GET,POST`, `/comments`, Comments)
		p.Route(`GET,POST`, `/favorites`, Favorites)
		p.Route(`GET,POST`, `/likes`, Likes)
	})

}
