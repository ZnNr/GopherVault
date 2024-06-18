package cmdutil // это имя пакета, которое вы можете использовать для удобства организации кода

import (
	"github.com/spf13/cobra" // импорт необходимой библиотеки для работы с командами
)

func GetFlagsValues(cmd *cobra.Command) (userName, bank, number, cv, login, password, metadata string) {
	userName, _ = cmd.Flags().GetString("user")
	bank, _ = cmd.Flags().GetString("bank")
	number, _ = cmd.Flags().GetString("number")
	cv, _ = cmd.Flags().GetString("cv")
	password, _ = cmd.Flags().GetString("password")
	metadata, _ = cmd.Flags().GetString("metadata")
	return
}
