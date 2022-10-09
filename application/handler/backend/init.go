package backend

import (
	_ "github.com/admpub/webx/application/handler/backend/official"

	// - register handler

	_ "github.com/admpub/webx/application/handler/backend/official/comment"
	_ "github.com/admpub/webx/application/handler/backend/official/manager"
	_ "github.com/admpub/webx/application/handler/backend/official/tool"
	_ "github.com/admpub/webx/application/handler/backend/official/user"
)
