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

// registerCmd представляет команду регистрации пользователей.
var registerCmd = &cobra.Command{
	Use:     "register",
	Short:   "Register in the GopherVault system.",
	Long:    `Register in the GopherVault system with provided login and password`,
	Example: "GopherVault register --login <user-system-login> --password <user-system-password>",
	Run:     registerHandler,
}

func registerHandler(cmd *cobra.Command, args []string) {
	// Загружаем переменные окружения из файла .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("error while getting envs: %s", err)
	}

	// Загружаем конфигурацию из переменных окружения
	var cfg models.Params
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("error while loading envs: %s\n", err)
	}

	// Получаем логин и пароль из флагов команды
	login, _ := cmd.Flags().GetString("login")
	password, _ := cmd.Flags().GetString("password")
	userCreds := models.User{
		Login:    login,
		Password: password,
	}

	// Преобразуем структуру пользователя в JSON
	body, err := json.Marshal(userCreds)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Отправляем POST запрос для регистрации пользователя
	resp, err := resty.New().R().
		SetHeader("Content-type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("http://%s:%s/auth/register", cfg.ApplicationHost, cfg.ApplicationPort))
	if err != nil {
		log.Printf(err.Error())
	}

	// Проверяем статус код ответа
	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code is not OK: %s\n", resp.Status())
		fmt.Println(resp.String())
		return
	}

	// Выводим сообщение об успешной регистрации пользователя
	fmt.Printf("user %q was successfully registered in goph-keeper", login)
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().String("login", "", "user login")
	registerCmd.Flags().String("password", "", "user password")
}
