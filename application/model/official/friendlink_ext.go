package official

import "github.com/admpub/webx/application/dbschema"

type FriendlinkExt struct {
	*dbschema.OfficialCommonFriendlink
	Category *dbschema.OfficialCommonCategory `db:"-,relation=id:category_id|gtZero"`
}
