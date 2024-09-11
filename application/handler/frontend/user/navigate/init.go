package navigate

import "github.com/coscms/webcore/registry/navigate"

func init() {
	navigate.Default.Frontend.Add(navigate.Left, LeftNavigate)
	navigate.Default.Frontend.Add(navigate.Top, TopNavigate)
}
