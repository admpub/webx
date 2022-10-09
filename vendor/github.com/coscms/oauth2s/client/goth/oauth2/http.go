package oauth2

import (
	"time"

	"github.com/webx-top/com"
)

var DefaultClient = com.HTTPClientWithTimeout(10 * time.Second)
