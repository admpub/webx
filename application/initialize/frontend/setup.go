package frontend

import (
	"fmt"

	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler/setup"
)

var MinCustomerID = 1000

func init() {
	setup.OnInstalled(func(ctx echo.Context) error {
		sqlStr := fmt.Sprintf("ALTER TABLE `official_customer` AUTO_INCREMENT=%d", MinCustomerID)
		_, err := factory.NewParam().DB().ExecContext(ctx, sqlStr)
		return err
	})
}
