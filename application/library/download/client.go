package download

import (
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/com/httpClientOptions"
)

var (
	Client = com.HTTPClientWithTimeout(30*time.Second, httpClientOptions.InsecureSkipVerify())
)
