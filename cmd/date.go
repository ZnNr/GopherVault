package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// buildDate представляет дату сборки приложения.
var buildDate = ""

// dateCmd представляет команду date.
var dateCmd = &cobra.Command{
	Use:     "build-date",
	Short:   "Показать дату сборки",
	Example: "GopherVault build-date",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Дата сборки: %s\n", buildDate)
	},
}

// init добавляет команду dateCmd к rootCmd
func init() {
	rootCmd.AddCommand(dateCmd)
}
