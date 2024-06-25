package cmd

import (
	"fmt"
	cmdutil "github.com/ZnNr/GopherVault/cmdutils"
	"github.com/ZnNr/GopherVault/internal/models"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// getCardCmd представляет команду getCard.
var getCardCmd = &cobra.Command{
	Use:     "get-card",
	Short:   "Get card info from GopherVault storage",
	Example: "GopherVault get-card --user <user-name> --number <card number>",
	Run:     getCardHandler,
}

func getCardHandler(cmd *cobra.Command, args []string) {
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

	resp, err := cmdutil.ExecutePostRequest(fmt.Sprintf("http://%s:%s/get/card", cfg.ApplicationHost, cfg.ApplicationPort), body)
	if err != nil {
		log.Printf(err.Error())
	}

	cmdutil.HandleResponse(resp, http.StatusOK)
}

func init() {
	rootCmd.AddCommand(getCardCmd)
	getCardCmd.Flags().String("user", "", "user name")
	getCardCmd.Flags().String("bank", "", "bank")
	getCardCmd.Flags().String("number", "", "number")
	getCardCmd.MarkFlagRequired("user")
}
