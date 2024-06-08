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

// addCredentialsCmd представляет команду add-credentials
var addCredentialsCmd = &cobra.Command{
	Use:   "add-credentials",
	Short: "Add a pair of login/password to GopherVault.",
	Long: `Add a pair of login/password to GopherVault database for
long-term storage. Only authorized users can use this command. The password is stored in the database in encrypted form.`,
	Example: "GopherVault add-credentials --user <user-name> --login <user-login> --password <password to store> --metadata <some description>",
	Run:     addCredentialsHandler,
}

// addCredentialsHandler обработчик команды добавления add-credentials
func addCredentialsHandler(cmd *cobra.Command, args []string) {
	// Загружаем переменные среды из файла .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("error while getting envs: %s", err)
	}

	var cfg models.Params
	// Загружаем переменные среды в структуру cfg
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("error while loading envs: %s\n", err)
	}

	// Получаем значения флагов команды
	userName, _ := cmd.Flags().GetString("user")
	login, _ := cmd.Flags().GetString("login")
	password, _ := cmd.Flags().GetString("password")
	metadata, _ := cmd.Flags().GetString("metadata")

	// Создаем объект модели Credentials
	requestCredentials := models.Credentials{
		UserName: userName,
		Login:    &login,
		Password: &password,
	}
	// Добавляем метаданные, если они указаны
	if metadata != "" {
		requestCredentials.Metadata = &metadata
	}

	// Преобразуем объект Credentials в JSON
	body, err := json.Marshal(requestCredentials)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Выполняем POST-запрос для сохранения учетных данных
	resp, err := resty.New().R().
		SetHeader("Content-type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("http://%s:%s/save/credentials", cfg.ApplicationHost, cfg.ApplicationPort))
	if err != nil {
		log.Printf(err.Error())
	}

	// Проверяем статус ответа
	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code is not OK: %s\n", resp.Status())
	}
	fmt.Println(resp.String())
}

func init() {
	rootCmd.AddCommand(addCredentialsCmd)
	addCredentialsCmd.Flags().String("user", "", "user name")
	addCredentialsCmd.Flags().String("login", "", "user login")
	addCredentialsCmd.Flags().String("password", "", "user password")
	addCredentialsCmd.Flags().String("metadata", "", "metadata")
	addCredentialsCmd.MarkFlagRequired("user")
	addCredentialsCmd.MarkFlagRequired("login")
	addCredentialsCmd.MarkFlagRequired("password")
}
