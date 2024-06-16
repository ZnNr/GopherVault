package cmdutil

import (
	"encoding/json"
	"github.com/ZnNr/GopherVault/internal/models"
	"log"
)

func ConvertToJSON(requestCard models.Card) []byte {
	body, err := json.Marshal(requestCard)
	if err != nil {
		log.Fatalf("ошибка при маршалинге запроса: %s", err.Error())
	}
	return body
}
