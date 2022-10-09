package customer

import "github.com/webx-top/db/lib/sqlbuilder"

var (
	PrivateFields = []string{
		`session_id`,
		`id_card_no`,
		`real_name`,
		`password`,
		`salt`,
		`safe_pwd`,
		`created`,
		`updated`,
		`file_num`,
		`file_size`,
	}

	PrivateFieldsWithMobileEmail = append([]string{
		`mobile_bind`,
		`mobile`,
		`email_bind`,
		`email`,
	}, PrivateFields...)

	//CusomterSafeFields 可公开的客户信息字段
	CusomterSafeFields = []interface{}{
		`id`,
		`name`,
		`gender`,
		`uid`,
		`group_id`,
		`avatar`,
		`online`,
		`following`,
		`followers`,
		`agent_level`,
	}

	CusomterSafeFieldsSelector = func(sel sqlbuilder.Selector) sqlbuilder.Selector {
		return sel.Columns(CusomterSafeFields...)
	}

	//UserSafeFields 可公开的后台用户信息字段
	UserSafeFields = []interface{}{
		`id`,
		`username`,
		`gender`,
		`avatar`,
		`online`,
	}

	UserSafeFieldsSelector = func(sel sqlbuilder.Selector) sqlbuilder.Selector {
		return sel.Columns(UserSafeFields...)
	}
)
