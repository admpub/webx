package navigate

import "github.com/admpub/nging/v4/application/registry/navigate"

func init() {
	navigate.Default.Frontend.Add(navigate.Left, LeftNavigate)
	navigate.Default.Frontend.Add(navigate.Top, TopNavigate)
}
