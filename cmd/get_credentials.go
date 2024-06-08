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

// getCredentialsCmd представляет команду get-credentials
var getCredentialsCmd = &cobra.Command{
	Use:   "get-credentials",
	Short: "Get a pair of login/password for specified user",
	Long: `Get a pair of login/password for specified user from GopherVault storage. 
Only authorized users can use this command`,
	Example: "GopherVault get-credentials --user user_name",
	Run:     getCredentialsHandler,
}

func getCredentialsHandler(cmd *cobra.Command, args []string) {
	// Загрузка переменных окружения
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Ошибка при загрузке переменных окружения: %s", err)
	}

	var cfg models.Params
	// Загрузка настроек из переменных окружения
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Ошибка при загрузке переменных окружения: %s\n", err)
	}

	userName, _ := cmd.Flags().GetString("user")
	userLogin, _ := cmd.Flags().GetString("login")
	requestUserCredentials := models.Credentials{
		UserName: userName,
	}
	if userLogin != "" {
		requestUserCredentials.Login = &userLogin
	}
	body, err := json.Marshal(requestUserCredentials)
	if err != nil {
		log.Fatalln(err.Error())
	}

	resp, err := resty.New().R().
		SetHeader("Content-type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("http://%s:%s/get/credentials", cfg.ApplicationHost, cfg.ApplicationPort))
	if err != nil {
		log.Printf(err.Error())
	}
	if resp.StatusCode() != http.StatusOK {
		log.Printf("Статус ответа не ОК: %s\n", resp.Status())
	}
	log.Printf(resp.String())

}

func init() {
	rootCmd.AddCommand(getCredentialsCmd)
	getCredentialsCmd.Flags().String("user", "", "user name")
	getCredentialsCmd.Flags().String("login", "", "user login")
	getCredentialsCmd.MarkFlagRequired("user")
}
