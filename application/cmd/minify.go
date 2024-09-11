package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/admpub/imageproxy"
	"github.com/coscms/webcore/cmd"
)

var minifyCmd = &cobra.Command{
	Use:   "minify",
	Short: "minify file",
	Long:  `Usage ./webx minify src.jpg dest.jpg`,
	RunE:  minifyRunE,
}

var minifyIMGOptions = imageproxy.Options{
	Quality: 70,
}

func minifyRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return cmd.Usage()
	}
	src, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	thumb, err := imageproxy.Transform(src, minifyIMGOptions)
	if err != nil {
		return err
	}
	err = os.WriteFile(args[1], thumb, 0644)
	return err
}

func init() {
	cmd.Add(minifyCmd)
}
