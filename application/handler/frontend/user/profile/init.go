package profile

import (
	"github.com/webx-top/echo"

	_ "github.com/admpub/webx/application/handler/frontend/user/wallet"
	"github.com/coscms/webfront/initialize/frontend"
)

func init() {
	frontend.RegisterToGroup(`/user`, func(u echo.RouteRegister) {
		p := u.Group(`/profile`)
		p.Route(`GET,POST`, `/detail`, Profile).SetName(`user.profile`)
		p.Route(`GET,POST`, `/settings`, Settings).SetName(`user.profile.settings`)
		p.Route(`GET,POST`, `/password`, Password).SetName(`user.profile.password`)
		p.Route(`GET,POST`, `/partial_authentication`, _Authentication)
		p.Route(`GET,POST`, `/binding`, Binding).SetName(`user.profile.binding`)
		p.Route(`POST`, `/follow`, Follow).SetName(`user.profile.follow`)
		p.Route(`POST`, `/unfollow`, Unfollow).SetName(`user.profile.unfollow`)
		p.Route(`GET,POST`, `/is_followed`, IsFollowed).SetName(`user.profile.is_followed`)
		p.Route(`GET,POST`, `/following`, Following).SetName(`user.profile.following`)
		p.Route(`GET,POST`, `/followers`, Followers).SetName(`user.profile.followers`)
		p.Route(`GET,POST`, `/comments`, Comments).SetName(`user.profile.comments`)
		p.Route(`GET,POST`, `/favorites`, Favorites).SetName(`user.profile.favorites`)
		p.Route(`GET,POST`, `/likes`, Likes).SetName(`user.profile.likes`)
	})

}
