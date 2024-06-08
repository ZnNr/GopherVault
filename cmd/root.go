package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// version представляет текущую версию приложения.
var version = "v0.0.0"

// rootCmd представляет основную команду при вызове без подкоманд.
var rootCmd = &cobra.Command{
	Use:     "GopherVault",
	Version: version,
	Short:   "GopherVault is a client-server system that allows the user to safely and securely store logins, passwords, binary data and other sensitive information.",
}

// Execute добавляет все дочерние команды к корневой команде и устанавливает флаги соответствующим образом.
// Эта функция вызывается из main.main(). Она должна быть вызвана только один раз для rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
