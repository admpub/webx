package cmd

import (
	"github.com/spf13/cobra"
	"github.com/webx-top/db/lib/factory"

	"github.com/admpub/nging/v4/application/cmd"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/webx/application/cmd/maker"
)

var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "Generate code",
	Long:  `Usage ./webx  make --group "official/b2c" --switchableFields "online"`,
	RunE:  makeRunE,
}

func makeRunE(cmd *cobra.Command, args []string) error {
	err := config.ParseConfig()
	if err != nil {
		return err
	}
	return maker.Make(maker.DefaultCLIConfig)
}

func init() {
	cmd.Add(makeCmd)
	makeCmd.Flags().StringVar(&maker.DefaultCLIConfig.Tables, "tables", "", "指定表名称,多个用半角逗号隔开。也可以备注该表的中文名称,英文名称与中文名称之间用半角冒号隔开,例如：actor:演员,area:区域")
	makeCmd.Flags().StringVar(&maker.DefaultCLIConfig.Group, "group", "official/film", "组名称,例如official/film,系统会自动将组名称中的“/”替换为“_”后作为表前缀")
	makeCmd.Flags().StringVar(&maker.DefaultCLIConfig.SwitchableFields, "switchableFields", "", "可切换状态的字段列表(各个字段用“,”分隔)")
	makeCmd.Flags().StringVar(&maker.DefaultCLIConfig.DBKey, "dbKey", factory.DefaultDBKey, "")
	makeCmd.Flags().BoolVar(&maker.DefaultCLIConfig.MustHasPrimaryKey, "mustHasPrimaryKey", false, "是否必须包含主键")
}
