package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
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
	// Загрузка переменных окружения
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Ошибка при загрузке окружения: %s", err)
	}

	var cfg models.Params
	// Загрузка настроек из переменных окружения
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
		Post(fmt.Sprintf("http://%s:%s/delete/card", cfg.ApplicationHost, cfg.ApplicationPort))
	if err != nil {
		log.Printf(err.Error())
	}
	if resp.StatusCode() != http.StatusOK {
		log.Printf("Статус ответа не ОК: %s\n", resp.Status())
	}
	log.Printf(resp.String())

}

func init() {
	rootCmd.AddCommand(deleteCardCmd)
	deleteCardCmd.Flags().String("user", "", "user name")
	deleteCardCmd.Flags().String("bank", "", "bank")
	deleteCardCmd.Flags().String("number", "", "card number")
	deleteCardCmd.MarkFlagRequired("user")
}
