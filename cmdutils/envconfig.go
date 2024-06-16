package cmdutil

import (
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

func LoadEnvVariables() models.Params {
	var cfg models.Params
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("ошибка при загрузке переменных окружения из файла: %s", err)
	}
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("ошибка при загрузке переменных окружения: %s", err)
	}
	return cfg
}
