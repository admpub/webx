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
	ReplyTarget *CommentAndExtraLite `db:"-,relation=id:reply_comment_id|gtZero,columns=owner_id:uint64&owner_type&content&contype&created:uint&updated:uint&display" json:",omitempty"`
}

type CommentLite struct {
	OwnerId   uint64 `db:"owner_id" bson:"owner_id" comment:"评论者ID" json:"owner_id" xml:"owner_id"`
	OwnerType string `db:"owner_type" bson:"owner_type" comment:"评论者类型(customer-前台客户;user-后台用户)" json:"owner_type" xml:"owner_type"`
	Content   string `db:"content" bson:"content" comment:"评论内容" json:"content" xml:"content"`
	Contype   string `db:"contype" bson:"contype" comment:"内容类型" json:"contype" xml:"contype"`
	Created   uint   `db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
	Updated   uint   `db:"updated" bson:"updated" comment:"编辑时间" json:"updated" xml:"updated"`
	Display   string `db:"display" bson:"display" comment:"显示" json:"display" xml:"display"`
}

type CommentAndExtraLite struct {
	*CommentLite
	User     *modelAuthor.User     `db:"-,relation=id:owner_id|gtZero|eq(owner_type:user),columns=id&username&avatar" json:",omitempty"`
	Customer *modelAuthor.Customer `db:"-,relation=id:owner_id|gtZero|eq(owner_type:customer),columns=id&name&avatar" json:",omitempty"`
}

func (c *CommentAndExtraLite) Name_() string {
	return dbschema.WithPrefix(`official_common_comment`)
}

type CommentAndExtra struct {
	*CommentAndReplyTarget
	FloorNumber      int
	ReplyFloorNumber int
	Extra            echo.H
}
