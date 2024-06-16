package cmdutil // это имя пакета, которое вы можете использовать для удобства организации кода

import (
	"github.com/spf13/cobra" // импорт необходимой библиотеки для работы с командами
	"log"
)

func getStringFlagValue(cmd *cobra.Command, flagName string) (string, error) {
	value, err := cmd.Flags().GetString(flagName)
	if err != nil {
		return "", err
	}
	return value, nil
}

// Другие вспомогательные функции для обработки флагов могут быть добавлены здесь

// Пример использования функции getStringFlagValue для обработки флагов
func ProcessFlags(cmd *cobra.Command) {
	userName, err := getStringFlagValue(cmd, "user")
	if err != nil {
		log.Printf("ошибка при получении значения флага user: %s", err.Error())
	}

	title, err := getStringFlagValue(cmd, "title")
	if err != nil {
		log.Printf("ошибка при получении значения флага title: %s", err.Error())
	}

	content, err := getStringFlagValue(cmd, "content")
	if err != nil {
		log.Printf("ошибка при получении значения флага content: %s", err.Error())
	}

	login, err := getStringFlagValue(cmd, "login")
	if err != nil {
		log.Printf("ошибка при получении значения флага login: %s", err.Error())
	}

	password, err := getStringFlagValue(cmd, "passworde")
	if err != nil {
		log.Printf("ошибка при получении значения флага password: %s", err.Error())
	}

	metadata, err := getStringFlagValue(cmd, "metadata")
	if err != nil {
		log.Printf("ошибка при получении значения флага metadata: %s", err.Error())
	}

	// Дополнительная обработка значений флагов здесь
	return

}

func GetFlagsValues(cmd *cobra.Command) (userName, bank, number, cv, login, password, metadata string) {
	userName, _ = cmd.Flags().GetString("user")
	bank, _ = cmd.Flags().GetString("bank")
	number, _ = cmd.Flags().GetString("number")
	cv, _ = cmd.Flags().GetString("cv")
	password, _ = cmd.Flags().GetString("password")
	metadata, _ = cmd.Flags().GetString("metadata")
	return
}
