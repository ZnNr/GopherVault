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
	"strings"

	"github.com/spf13/cobra"
)

// addCardCmd представляет команду add-card
var addCardCmd = &cobra.Command{
	Use:   "add-card",
	Short: "Добавление информации о банковской карте в GopherVault.",
	Long: `Добавляет информацию о банковской карте (название банка, номер карты, CV, пароль и метаданные) в базу данных goph-keeper для
долгосрочного хранения. Только авторизованные пользователи могут использовать эту команду. Пароль и CV хранятся в базе данных в зашифрованном виде.`,
	Example: "gGopherVault add-card --user user-name --bank alpha --number 1111222233334444 --cv 123 --password 1243",
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("ошибка при загрузке переменных окружения из файла: %s", err)
		}

		var cfg models.Params
		if err := envconfig.Process("", &cfg); err != nil {
			log.Fatalf("ошибка при обработке переменных окружения: %s\n", err)
		}

		// Получение значений флагов из командной строки
		userName, _ := cmd.Flags().GetString("user")
		bank, _ := cmd.Flags().GetString("bank")
		number, _ := cmd.Flags().GetString("number")
		cv, _ := cmd.Flags().GetString("cv")
		password, _ := cmd.Flags().GetString("password")
		metadata, _ := cmd.Flags().GetString("metadata")

		// Проверка наличия всех обязательных значений
		if strings.TrimSpace(userName) == "" || strings.TrimSpace(bank) == "" || strings.TrimSpace(number) == "" || strings.TrimSpace(cv) == "" || strings.TrimSpace(password) == "" {
			log.Fatalln("имя пользователя, название банка, номер карты, CV и пароль не должны быть пустыми")
		}
		if len(number) != 16 {
			log.Fatalln("идентификационный номер пластиковой карты должен состоять из 16 цифр.")
		}
		if len(cv) != 3 {
			log.Fatalln("CV-код пластиковой карты должен состоять из 3 цифр.")
		}

		// Создание структуры запроса для карты
		requestCard := models.Card{
			UserName: userName,
			BankName: &bank,
			Number:   &number,
			CV:       &cv,
			Password: &password,
		}
		if metadata != "" {
			requestCard.Metadata = &metadata
		}

		// Преобразование в JSON и отправка запроса на сервер
		body, err := json.Marshal(requestCard)
		if err != nil {
			log.Fatalf("ошибка при маршалинге запроса: %s", err.Error())
		}

		resp, err := resty.New().R().
			SetHeader("Content-type", "application/json").
			SetBody(body).
			Post(fmt.Sprintf("http://%s:%s/save/card", cfg.ApplicationHost, cfg.ApplicationPort))
		if err != nil {
			log.Printf("ошибка при выполнении запроса: %s", err.Error())
		}

		if resp.StatusCode() != http.StatusOK {
			log.Printf("код состояния не ОК: %s\n", resp.Status())
		}

		fmt.Println(resp.String())
	},
}

func init() {
	rootCmd.AddCommand(addCardCmd)

	// Определение флагов и их обязательность
	addCardCmd.Flags().String("user", "", "имя пользователя")
	addCardCmd.Flags().String("bank", "", "банк")
	addCardCmd.Flags().String("number", "", "номер карты")
	addCardCmd.Flags().String("cv", "", "CV карты")
	addCardCmd.Flags().String("password", "", "пароль карты")
	addCardCmd.Flags().String("metadata", "", "метаданные")
	addCardCmd.MarkFlagRequired("user")
	addCardCmd.MarkFlagRequired("bank")
	addCardCmd.MarkFlagRequired("number")
	addCardCmd.MarkFlagRequired("cv")
	addCardCmd.MarkFlagRequired("password")
}
