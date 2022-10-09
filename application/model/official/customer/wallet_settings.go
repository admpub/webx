package customer

import "github.com/webx-top/echo"

func NewWalletSettings() *WalletSettings {
	return &WalletSettings{}
}

type WalletSettings struct {
	On            bool
	MinAmount     float64
	DefaultAmount float64
}

func (s *WalletSettings) FromStore(r echo.H) *WalletSettings {
	s.On = r.Bool("On", false)
	s.MinAmount = r.Float64("MinAmount")
	s.DefaultAmount = r.Float64("DefaultAmount")
	return s
}
