package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ZnNr/GopherVault/internal/models"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/spf13/cobra"
)

// getCardCmd представляет команду getCard.
var getCardCmd = &cobra.Command{
	Use:     "get-card",
	Short:   "Get card info from GopherVault storage",
	Example: "GopherVault get-card --user <user-name> --number <card number>",
	Run: func(cmd *cobra.Command, args []string) {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("Ошибка при загрузке окружения: %s", err)
		}

		var cfg models.Params
		if err := envconfig.Process("", &cfg); err != nil {
			log.Fatalf("Ошибка при загрузке окружения: %s\n", err)
		}

		userName, _ := cmd.Flags().GetString("user")
		bank, _ := cmd.Flags().GetString("bank")
		number, _ := cmd.Flags().GetString("number")
		requestCard := models.Card{
			UserName: userName,
		}
		if bank != "" {
			requestCard.BankName = &bank
		}
		if number != "" {
			requestCard.Number = &number
		}
		body, err := json.Marshal(requestCard)
		if err != nil {
			log.Fatalln(err.Error())
		}

		resp, err := resty.New().R().
			SetHeader("Content-type", "application/json").
			SetBody(body).
			Post(fmt.Sprintf("http://%s:%s/get/card", cfg.ApplicationHost, cfg.ApplicationPort))
		if err != nil {
			log.Printf(err.Error())
		}
		if resp.StatusCode() != http.StatusOK {
			log.Printf("Статус ответа не OK: %s\n", resp.Status())
		}
		log.Printf(resp.String())
	},
}

func init() {
	rootCmd.AddCommand(getCardCmd)
	getCardCmd.Flags().String("user", "", "имя пользователя")
	getCardCmd.Flags().String("bank", "", "банк")
	getCardCmd.Flags().String("number", "", "номер карты")
	getCardCmd.MarkFlagRequired("user")
}
