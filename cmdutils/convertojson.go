package cmdutil

import (
	"encoding/json"
	"github.com/ZnNr/GopherVault/internal/models"
	"log"
)

func ConvertToJSONRequestCards(requestCard models.Card) []byte {
	body, err := json.Marshal(requestCard)
	if err != nil {
		log.Fatalf("ошибка при маршалинге запроса: %s", err.Error())
	}
	return body
}

func ConvertToJSONRequestNotes(requestNote models.Note) []byte {
	body, err := json.Marshal(requestNote)
	if err != nil {
		log.Fatalf("ошибка при маршалинге запроса: %s", err.Error())
	}
	return body
}

func ConvertToJSONRequestCredential(userCreds models.Credentials) []byte {
	body, err := json.Marshal(userCreds)
	if err != nil {
		log.Fatalf("ошибка при маршалинге запроса: %s", err.Error())
	}
	return body
}

func ConvertToJSONRequestUserCredential(requestCredential models.User) []byte {
	body, err := json.Marshal(requestCredential)
	if err != nil {
		log.Fatalf("ошибка при маршалинге запроса: %s", err.Error())
	}
	return body
}
