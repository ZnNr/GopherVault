package cmd

import (
	"fmt"
	"github.com/ZnNr/GopherVault/cmdutils"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"strings"
)

// Определение команды add-card
var addCardCmd = &cobra.Command{
	Use:   "add-card",
	Short: "Add bank card info to GopherVault.",
	Long: `Add bank card info (bank name, card number, cv, password and metadata) to GopherVault database for
long-term storage. Only authorized users can use this command. Password and cv are stored in the database in the encrypted form.`,
	Example: "GopherVault add-card --user user-name --bank alpha --number 1111222233334444 --cv 123 --password 1243",
	Run:     addCardHandler,
}

// Обработчик команды add-card
func addCardHandler(cmd *cobra.Command, args []string) {
	cmdutil.LoadEnvVariables()

	// Получение значений флагов из командной строки
	userName, bank, number, cv, password, metadata, _ := cmdutil.GetFlagsValues(cmd)

	// Проверка наличия всех обязательных значений
	checkRequiredValues(userName, bank, number, cv, password)

	// Создание структуры запроса для карты
	requestCard := createCardRequest(userName, bank, number, cv, password, metadata)

	// Преобразование в JSON и отправка запроса на сервер
	body := cmdutil.ConvertToJSONRequestCards(requestCard)
	sendPostRequest(body)
}

func sendPostRequest(body []byte) {
	var cfg models.Params
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("ошибка при обработке переменных окружения: %s\n", err)
	}

	resp, err := resty.New().R().
		SetHeader("Content-type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("http://%s:%s/save/card", cfg.ApplicationHost, cfg.ApplicationPort))
	if err != nil {
		log.Printf("ошибка при выполнении запроса: %s", err.Error())
	}

	if resp != nil && resp.StatusCode() != http.StatusOK {
		log.Printf("код состояния не ОК: %s\n", resp.Status())
	}

	log.Println(resp.String())
}

func checkRequiredValues(userName, bank, number, cv, password string) {
	if strings.TrimSpace(userName) == "" || strings.TrimSpace(bank) == "" || strings.TrimSpace(number) == "" || strings.TrimSpace(cv) == "" || strings.TrimSpace(password) == "" {
		log.Fatalln("имя пользователя, название банка, номер карты, CV и пароль не должны быть пустыми")
	}
	if len(number) != 16 {
		log.Fatalln("идентификационный номер пластиковой карты должен состоять из 16 цифр.")
	}
	if len(cv) != 3 {
		log.Fatalln("CV-код пластиковой карты должен состоять из 3 цифр.")
	}
}

func createCardRequest(userName, bank, number, cv, password, metadata string) models.Card {
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
	return requestCard
}

func init() {
	rootCmd.AddCommand(addCardCmd)

	// Определение флагов и их обязательность
	addCardCmd.Flags().String("user", "", "user name")
	addCardCmd.Flags().String("bank", "", "bank")
	addCardCmd.Flags().String("number", "", "card number")
	addCardCmd.Flags().String("cv", "", "card cv")
	addCardCmd.Flags().String("password", "", "card password")
	addCardCmd.Flags().String("metadata", "", "metadata")
	addCardCmd.MarkFlagRequired("user")
	addCardCmd.MarkFlagRequired("bank")
	addCardCmd.MarkFlagRequired("number")
	addCardCmd.MarkFlagRequired("cv")
	addCardCmd.MarkFlagRequired("password")
}
