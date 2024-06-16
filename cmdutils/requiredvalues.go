package cmdutil

import (
	"log"
	"strings"
)

func CheckRequiredValues(userName, bank, number, cv, password string) {
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
