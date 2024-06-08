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

// updateCredentialsCmd представляет команду updateCredentials
var updateCredentialsCmd = &cobra.Command{
	Use:     "update-credentials",
	Short:   "Update user credentials for the provided login.",
	Example: "GopherVault update-credentials --user <user-name> --login <saved-login> --password <new-password>",
	Run:     updateCredentialsHandler,
}

// updateCredentialsHandler обработчик команды обновления учетных данных
func updateCredentialsHandler(cmd *cobra.Command, args []string) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Ошибка при загрузке переменных окружения: %s", err)
	}

	var cfg models.Params
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %s\n", err)
	}

	userName, _ := cmd.Flags().GetString("user")
	login, _ := cmd.Flags().GetString("login")
	password, _ := cmd.Flags().GetString("password")
	metadata, _ := cmd.Flags().GetString("metadata")

	requestCredentials := models.Credentials{
		UserName: userName,
		Login:    &login,
		Password: &password,
		Metadata: &metadata,
	}

	body, err := json.Marshal(requestCredentials)
	if err != nil {
		log.Fatalf(err.Error())
	}

	resp, err := sendUpdateCredentialsRequest(cfg, body)
	if err != nil {
		log.Printf(err.Error())
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("Статус ответа не 'OK': %s\n", resp.Status())
	}

	log.Println(resp.String())
}

// sendUpdateCredentialsRequest отправляет запрос на обновление учетных данных
func sendUpdateCredentialsRequest(cfg models.Params, body []byte) (*resty.Response, error) {
	return resty.New().R().
		SetHeader("Content-type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("http://%s:%s/update/credentials", cfg.ApplicationHost, cfg.ApplicationPort))
}

func init() {
	rootCmd.AddCommand(updateCredentialsCmd)
	updateCredentialsCmd.Flags().String("user", "", "user name")
	updateCredentialsCmd.Flags().String("login", "", "user login")
	updateCredentialsCmd.Flags().String("password", "", "user password")
	updateCredentialsCmd.Flags().String("metadata", "", "metadata")
	updateCredentialsCmd.MarkFlagRequired("user")
	updateCredentialsCmd.MarkFlagRequired("login")
	updateCredentialsCmd.MarkFlagRequired("password")
}
