package user

import (
	"testing"
	"time"

	"github.com/coscms/webcore/library/testutils"
	"github.com/stretchr/testify/require"
	"github.com/webx-top/echo/defaults"
)

func TestQRSignIn(t *testing.T) {
	testutils.InitConfig()
	qsi := QRSignIn{
		SessionID:     `477V7JN2QGSFTFHJ6JHBWZ6M7H7OX7KVEL4TPBODYRDED4GIZM5A`,
		SessionMaxAge: 86400 * 365,
		Expires:       time.Now().Add(time.Minute * 10).Unix(),
		IPAddress:     `127.0.0.1`,
		UserAgent:     `Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:133.0) Gecko/20100101 Firefox/133.0`,
		Platform:      `windows`,
		Scense:        `qrcode_web`,
		DeviceNo:      `1323232324545656075970860458045045083949`,
	}
	r, err := qsi.Encode()
	require.NoError(t, err)
	t.Logf(`%s: %d`, r, len(r))

	qrcode := GenerateUniqueKey(qsi.IPAddress, qsi.UserAgent)
	t.Logf(`%s: %d`, qrcode, len(qrcode))

	// -------------------------------

	cs := GetQRSignInCase(`cache`)
	ctx := defaults.NewMockContext()
	qrcode, err = cs.Encode(ctx, qsi)
	require.NoError(t, err)
	t.Logf(`%s: %d`, qrcode, len(qrcode))
	qsi2, err := cs.Decode(ctx, qrcode)
	require.NoError(t, err)
	require.Equal(t, qsi, qsi2)
}
