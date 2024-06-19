package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ZnNr/GopherVault/internal/models"
	"sync"
	"time"

	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

var jwtKey = []byte("my_secret_key")

type handler struct {
	db        models.Storage
	log       *zap.SugaredLogger
	cookiesMu sync.Mutex
	cookies   map[string]string
}

func New(db models.Storage, log *zap.SugaredLogger) *handler {
	return &handler{
		db:        db,
		log:       log,
		cookiesMu: sync.Mutex{},
		cookies:   make(map[string]string),
	}
}

// LoginHandler обрабатывает запросы на аутентификацию пользователей
func (h *handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	h.cookiesMu.Lock()
	defer h.cookiesMu.Unlock()

	// Используем контекст из запроса
	ctx := r.Context()

	// Извлекаем данные пользователя из тела запроса
	user, err := parseUserInput(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем пароль пользователя
	if err = h.db.Login(ctx, user.Login, user.Password); err != nil {
		message, status := handleUserError(user.Login, err)
		http.Error(w, message, status)
		return
	}

	// Создаем JWT токен для пользователя, добавляем заголовок Authorization и устанавливаем cookie
	expirationTime := time.Now().Add(time.Hour)
	token, err := createToken(user.Login, expirationTime)
	if err != nil {
		http.Error(w, fmt.Sprintf("ошибка при создании токена для пользователя: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})

	// Отправляем успешный статус ответа
	w.WriteHeader(http.StatusOK)

	// Сохраняем токен в памяти и записываем информацию в лог
	h.cookies[user.Login] = fmt.Sprintf("Bearer %s", token)
	h.log.Infof("пользователь %q успешно вошел в систему", user.Login)
}

// RegisterHandler обрабатывает запросы на регистрацию новых пользователей
func (h *handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	h.cookiesMu.Lock()
	defer h.cookiesMu.Unlock()

	// Используем контекст из запроса
	ctx := r.Context()

	// Извлекаем данные пользователя из тела запроса
	user, err := parseUserInput(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Регистрируем пользователя в системе
	if err = h.db.Register(ctx, user.Login, user.Password); err != nil {
		message, status := handleUserError(user.Login, err)
		http.Error(w, message, status)
		return
	}

	// Создаем JWT токен для пользователя, добавляем заголовок Authorization и устанавливаем cookie
	expirationTime := time.Now().Add(time.Hour)
	token, err := createToken(user.Login, expirationTime)
	if err != nil {
		http.Error(w, fmt.Sprintf("ошибка при создании токена для пользователя: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})

	// Отправляем успешный статус ответа
	w.WriteHeader(http.StatusOK)

	// Сохраняем токен в памяти и записываем информацию в лог
	h.cookies[user.Login] = fmt.Sprintf("Bearer %s", token)
	h.log.Infof("пользователь %q успешно зарегистрирован", user.Login)
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

// SaveCardHandler обрабатывает запросы на сохранение карточки пользователя
func (h *handler) SaveCardHandler(w http.ResponseWriter, r *http.Request) {
	h.cookiesMu.Lock()
	defer h.cookiesMu.Unlock()
	// Создаем контекст для запроса
	ctx := r.Context()

	// Читаем тело запроса
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Распаковываем данные из тела запроса в структуру Card
	var requestCard models.Card
	if err := json.Unmarshal(buf.Bytes(), &requestCard); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Сохраняем карточку пользователя в хранилище goph-keeper
	if err := h.db.SaveCard(ctx, requestCard); err != nil {
		message, status := handleUserError(requestCard.UserName, err)
		http.Error(w, message, status)
		return
	}

	// Формируем ответ
	response := fmt.Sprintf("Карточка для пользователя %q успешно сохранена", requestCard.UserName)

	// Отправляем ответ клиенту
	if _, err := io.WriteString(w, response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetCardHandler обрабатывает запросы на получение карточек пользователя
func (h *handler) GetCardHandler(w http.ResponseWriter, r *http.Request) {
	h.cookiesMu.Lock()
	defer h.cookiesMu.Unlock()

	// Создаем контекст для запроса
	ctx := r.Context()

	// Читаем тело запроса для получения имени пользователя
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Распаковываем данные из тела запроса в структуру Card
	var cardRequest models.Card
	if err := json.Unmarshal(body, &cardRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем карточки пользователя из хранилища goph-keeper
	cards, err := h.db.GetCard(ctx, cardRequest)
	if err != nil {
		message, status := handleUserError(cardRequest.UserName, err)
		http.Error(w, message, status)
		return
	}

	// Формируем ответ
	cardsResponse, err := json.Marshal(cards)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ клиенту
	if _, err := io.WriteString(w, string(cardsResponse)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteCardHandler обрабатывает запросы на удаление карточек пользователя
func (h *handler) DeleteCardHandler(w http.ResponseWriter, r *http.Request) {
	h.cookiesMu.Lock()
	defer h.cookiesMu.Unlock()

	// Создаем контекст для запроса
	ctx := r.Context()

	// Читаем тело запроса для получения информации о карточке
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Распаковываем данные из тела запроса в структуру Card
	var cardRequest models.Card
	if err = json.Unmarshal(body, &cardRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Удаляем карточки пользователя из хранилища goph-keeper
	if err = h.db.DeleteCards(ctx, cardRequest); err != nil {
		message, status := handleUserError(cardRequest.UserName, err)
		http.Error(w, message, status)
		return
	}

	// Формируем сообщение об успешном удалении карточек
	response := generateDeleteResponse(cardRequest)

	// Отправляем ответ клиенту
	if _, err = io.WriteString(w, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// generateDeleteResponse генерирует сообщение об успешном удалении карточек
func generateDeleteResponse(card models.Card) string {
	response := fmt.Sprintf("Cards for user %q were successfully deleted", card.UserName)
	if card.BankName != nil {
		response = fmt.Sprintf("Cards of %q bank for user %q were successfully deleted", *card.BankName, card.UserName)
	} else if card.Number != nil {
		response = fmt.Sprintf("Cards with number %q for user %q were successfully deleted", *card.Number, card.UserName)
	}
	return response
}
