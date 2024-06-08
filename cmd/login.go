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

// loginCmd представляет команду login
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the GopherVault system",
	Long: `Login to the GopherVault system with specified login and password. 
Only registered users can run this command`,
	Example: "GopherVault login --login <user-system-login> --password <user-system-password>`",
	Run:     loginHandler,
}

func loginHandler(cmd *cobra.Command, args []string) {
	// Загрузка переменных окружения
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Ошибка при загрузке переменных окружения: %s", err)
	}

	var cfg models.Params
	// Загрузка настроек из переменных окружения
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Ошибка при загрузке переменных окружения: %s\n", err)
	}

	login, _ := cmd.Flags().GetString("login")
	password, _ := cmd.Flags().GetString("password")
	userCreds := models.User{
		Login:    login,
		Password: password,
	}

	// Преобразование пользовательских данных в JSON
	body, err := json.Marshal(userCreds)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Отправка POST запроса для аутентификации пользователя
	resp, err := resty.New().R().
		SetHeader("Content-type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("http://%s:%s/auth/login", cfg.ApplicationHost, cfg.ApplicationPort))
	if err != nil {
		log.Printf(err.Error())
	}

	// Проверка статуса ответа
	if resp.StatusCode() != http.StatusOK {
		log.Printf("Статус ответа не является OK: %s\n", resp.Status())
		fmt.Println(resp.String())
		return
	}

	// Вывод сообщения об успешном входе пользователя
	fmt.Printf("Пользователь %q успешно вошел в систему GopherVault\n", login)
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().String("login", "", "user login")
	loginCmd.Flags().String("password", "", "user password")
	loginCmd.MarkFlagRequired("login")
	loginCmd.MarkFlagRequired("password")
}
