package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ZnNr/GopherVault/internal/models"

	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

var jwtKey = []byte("my_secret_key")

type handler struct {
	credentialsStorage models.CredentialsStorage
	noteStorage        models.NoteStorage
	cardStorage        models.CardStorage
	authService        models.AuthenticationService
	log                *zap.SugaredLogger
	cookies            map[string]string
}

func New(credentialsStorage models.CredentialsStorage, noteStorage models.NoteStorage, cardStorage models.CardStorage, authService models.AuthenticationService, log *zap.SugaredLogger) *handler {
	return &handler{
		credentialsStorage: credentialsStorage,
		noteStorage:        noteStorage,
		cardStorage:        cardStorage,
		authService:        authService,
		log:                log,
		cookies:            make(map[string]string),
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
}

func (h *handler) GetUserCredentials(w http.ResponseWriter, r *http.Request) {
}

func (h *handler) SaveUserCredentials(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) DeleteUserCredentials(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) UpdateUserCredentials(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) SaveUserNote(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) GetUserNote(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) DeleteUserNotes(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) UpdateUserNote(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) SaveCard(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) GetCard(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) DeleteCard(w http.ResponseWriter, r *http.Request) {

}

// CheckAuthorization проверяет авторизацию текущего пользователя.
func (h *handler) CheckAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Чтение и декодирование JSON из тела запроса в структуру Credentials
		var user models.Credentials
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Println("CheckAuthorization: error decoding JSON")
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Проверка наличия куки для пользователя
		if h.cookies[user.UserName] == "" {
			http.Error(w, fmt.Sprintf("User %q is not authorized", user.UserName), http.StatusUnauthorized)
			return
		}

		// Проверка токена
		token, err := extractJwtToken(h.cookies[user.UserName])
		if err != nil {
			message, status := handleUserError(user.UserName, err)
			http.Error(w, message, status)
			return
		}

		// Проверка валидности токена
		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Добавление заголовка Authorization и передача управления следующему обработчику
		w.Header().Add("Authorization", h.cookies[user.UserName])
		r.Body = io.NopCloser(bytes.NewBuffer([]byte{})) // сброс тела запроса для последующего чтения
		next.ServeHTTP(w, r)
	})
}
