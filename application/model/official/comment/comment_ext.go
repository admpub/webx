package comment

import (
	"github.com/admpub/webx/application/dbschema"
	"github.com/webx-top/echo"
)

type CommentAndCustomer struct {
	*dbschema.OfficialCommonComment
	Customer *dbschema.OfficialCustomer
}

type CommentAndExtra struct {
	*dbschema.OfficialCommonComment
	FloorNumber      int
	ReplyFloorNumber int
	Extra            echo.H
}
