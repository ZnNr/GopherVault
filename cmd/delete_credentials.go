package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

// deleteCredentialsCmd представляет команду deleteCredentials
var deleteCredentialsCmd = &cobra.Command{
	Use:     "delete-credentials",
	Short:   "Delete credentials for user from GopherVault storage",
	Example: "GopherVault delete-credentials --user <user-name> --login <user-login>",
	Run:     deleteCredentialsHandler,
}

func deleteCredentialsHandler(cmd *cobra.Command, args []string) {
	// Загружаем переменные среды из файла .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error while getting envs: %s", err)
	}

	var cfg models.Params
	// Загружаем переменные среды в структуру cfg
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Error while loading envs: %s\n", err)
	}

	userName, _ := cmd.Flags().GetString("user")
	login, _ := cmd.Flags().GetString("login")

	// Создаем объект модели Credentials для запроса
	requestUserCredentials := models.Credentials{
		UserName: userName,
	}
	// Добавляем логин, если он указан
	if login != "" {
		requestUserCredentials.Login = &login
	}

	// Преобразуем объект в json
	body, err := json.Marshal(requestUserCredentials)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(string(body))

	// Отправляем POST-запрос на удаление учетных данных
	resp, err := resty.New().R().
		SetHeader("Content-type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("http://%s:%s/delete/credentials", cfg.ApplicationHost, cfg.ApplicationPort))
	if err != nil {
		log.Printf(err.Error())
	}
	if resp.StatusCode() != http.StatusOK {
		log.Printf("Status code is not OK: %s\n", resp.Status())
	}
	log.Printf(resp.String())

}

func init() {
	rootCmd.AddCommand(deleteCredentialsCmd)
	deleteCredentialsCmd.Flags().String("user", "", "user name")
	deleteCredentialsCmd.Flags().String("login", "", "user login")
	deleteCredentialsCmd.MarkFlagRequired("user")
}
