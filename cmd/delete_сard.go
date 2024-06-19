package cmd

import (
	"fmt"
	cmdutil "github.com/ZnNr/GopherVault/cmdutils"
	"github.com/ZnNr/GopherVault/internal/models"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// deleteCardCmd представляет команду deleteCard.
var deleteCardCmd = &cobra.Command{
	Use:     "delete-card",
	Short:   "Delete card info from GopherVault storage",
	Example: "GopherVault  delete-card --user user-name --bank alpha",
	Run:     deleteCardHandler,
}

func deleteCardHandler(cmd *cobra.Command, args []string) {
	cfg := cmdutil.LoadEnvVariables()
	userName, bank, number, _, _, _, _, _ := cmdutil.GetFlagsValues(cmd)

	requestCard := models.Card{
		UserName: userName,
	}
	if bank != "" {
		requestCard.BankName = &bank
	}
	if number != "" {
		requestCard.Number = &number
	}
	body := cmdutil.ConvertToJSONRequestCards(requestCard)

	resp, err := cmdutil.ExecutePostRequest(fmt.Sprintf("http://%s:%s/delete/card", cfg.ApplicationHost, cfg.ApplicationPort), body)
	if err != nil {
		log.Printf(err.Error())
	}

	cmdutil.HandleResponse(resp, http.StatusOK)
}

func init() {
	rootCmd.AddCommand(deleteCardCmd)
	deleteCardCmd.Flags().String("user", "", "user name")
	deleteCardCmd.Flags().String("bank", "", "bank")
	deleteCardCmd.Flags().String("number", "", "card number")
	deleteCardCmd.MarkFlagRequired("user")
}
