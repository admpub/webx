package comment

import (
	"github.com/webx-top/echo"

	"github.com/admpub/webx/application/dbschema"
	modelAuthor "github.com/admpub/webx/application/model/author"
)

type CommentAndCustomer struct {
	*dbschema.OfficialCommonComment
	Customer *dbschema.OfficialCustomer
}

type CommentAndReplyTarget struct {
	*dbschema.OfficialCommonComment
	ReplyTarget *CommentAndExtraLite `db:"-,relation=id:reply_comment_id|gtZero" json:",omitempty"`
}

type CommentAndExtraLite struct {
	*dbschema.OfficialCommonComment
	User     *modelAuthor.User     `db:"-,relation=id:owner_id|gtZero|eq(owner_type:user),columns=id&username&avatar" json:",omitempty"`
	Customer *modelAuthor.Customer `db:"-,relation=id:owner_id|gtZero|eq(owner_type:customer),columns=id&name&avatar" json:",omitempty"`
}

func (c *CommentAndExtraLite) Name_() string {
	if c.OfficialCommonComment == nil {
		c.OfficialCommonComment = &dbschema.OfficialCommonComment{}
	}
	return c.OfficialCommonComment.Name_()
}

type CommentAndExtra struct {
	*CommentAndReplyTarget
	FloorNumber      int
	ReplyFloorNumber int
	Extra            echo.H
}
