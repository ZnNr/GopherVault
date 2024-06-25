package cmdutil

import (
	"github.com/spf13/cobra"
)

func GetFlagsValues(cmd *cobra.Command) (userName, bank, number, cv, login, password, сardType, metadata string) {
	userName, _ = cmd.Flags().GetString("user")
	bank, _ = cmd.Flags().GetString("bank")
	number, _ = cmd.Flags().GetString("number")
	cv, _ = cmd.Flags().GetString("cv")
	login, _ = cmd.Flags().GetString("login")
	password, _ = cmd.Flags().GetString("password")
	сardType, _ = cmd.Flags().GetString("type")
	metadata, _ = cmd.Flags().GetString("metadata")
	return
}
