package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ZnNr/GopherVault/internal/models"
	"time"

	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

var jwtKey = []byte("my_secret_key")

type handler struct {
	db      models.Storage
	log     *zap.SugaredLogger
	cookies map[string]string
}

func New(db models.Storage, log *zap.SugaredLogger) *handler {
	return &handler{
		db:      db,
		log:     log,
		cookies: make(map[string]string),
	}
}

// LoginHandler обрабатывает запросы на аутентификацию пользователей
func (h *handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Создаем контекст для запроса
	ctx := context.Background()

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
	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

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

// GetUserCredentialsHandler обрабатывает запросы на получение учетных данных пользователя
func (h *handler) GetUserCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Создаем контекст для запроса
	ctx := context.Background()

	// Извлекаем данные пользователя из тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userCredentialsRequest models.Credentials
	if err = json.Unmarshal(body, &userCredentialsRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем учетные данные пользователя из хранилища
	creds, err := h.db.GetCredentials(ctx, userCredentialsRequest)
	if err != nil {
		message, status := handleUserError(userCredentialsRequest.UserName, err)
		http.Error(w, message, status)
		return
	}

	// Формируем ответ
	userCredentialsJSON, err := json.Marshal(creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err = io.WriteString(w, string(userCredentialsJSON)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный статус ответа
	w.WriteHeader(http.StatusOK)
}

// SaveUserCredentialsHandler обрабатывает запросы на сохранение учетных данных пользователя
func (h *handler) SaveUserCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Создаем контекст для запроса
	ctx := context.Background()

	// Получаем данные из тела запроса
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Распаковываем полученные данные в структуру Credentials
	var requestCredentials models.Credentials
	if err := json.Unmarshal(buf.Bytes(), &requestCredentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем наличие логина и пароля
	if requestCredentials.Login == nil || requestCredentials.Password == nil {
		http.Error(w, "логин и пароль не должны быть пустыми", http.StatusBadRequest)
		return
	}

	// Сохраняем учетные данные пользователя в хранилище
	if err := h.db.SaveCredentials(ctx, requestCredentials); err != nil {
		message, status := handleUserError(requestCredentials.UserName, err)
		http.Error(w, message, status)
		return
	}

	// Возвращаем успешный ответ
	if _, err := io.WriteString(w, fmt.Sprintf("Учетные данные для пользователя %q сохранены", requestCredentials.UserName)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный статус ответа
	w.WriteHeader(http.StatusOK)
}

// DeleteUserCredentialsHandler обрабатывает запросы на удаление учетных данных пользователя
func (h *handler) DeleteUserCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Создаем контекст для запроса
	ctx := context.Background()

	// Читаем тело запроса для получения имени пользователя
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Распаковываем данные из тела запроса в структуру Credentials
	var userCredentialsRequest models.Credentials
	if err := json.Unmarshal(body, &userCredentialsRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Удаляем учетные данные из хранилища
	if err := h.db.DeleteCredentials(ctx, userCredentialsRequest); err != nil {
		message, status := handleUserError(userCredentialsRequest.UserName, err)
		http.Error(w, message, status)
		return
	}

	// Формируем сообщение об успешном удалении учетных данных
	var response string
	if userCredentialsRequest.Login != nil {
		response = fmt.Sprintf("Учетные данные для пользователя %q с логином %q были успешно удалены", userCredentialsRequest.UserName, *userCredentialsRequest.Login)
	} else {
		response = fmt.Sprintf("Учетные данные для пользователя %q были успешно удалены", userCredentialsRequest.UserName)
	}

	// Отправляем ответ клиенту
	if _, err := io.WriteString(w, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем успешный статус ответа
	w.WriteHeader(http.StatusOK)
}

// UpdateUserCredentialsHandler обрабатывает запросы на обновление учетных данных пользователя
func (h *handler) UpdateUserCredentialsHandler(w http.ResponseWriter, r *http.Request) {

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Создаем контекст для запроса
	ctx := context.Background()

	// Читаем тело запроса для получения данных обновления учетных данных пользователя
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Распаковываем данные из тела запроса в структуру Credentials
	var requestCredentials models.Credentials
	if err := json.Unmarshal(buf.Bytes(), &requestCredentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем обязательные поля (логин и пароль)
	if requestCredentials.Login == nil || requestCredentials.Password == nil {
		http.Error(w, "login and password should not be empty", http.StatusBadRequest)
		return
	}

	// Обновляем учетные данные пользователя в хранилище
	if err := h.db.UpdateCredentials(ctx, requestCredentials); err != nil {
		message, status := handleUserError(requestCredentials.UserName, err)
		http.Error(w, message, status)
		return
	}

	// Формируем сообщение об успешном обновлении учетных данных
	response := fmt.Sprintf("Учетные данные для пользователя %q были успешно обновлены", requestCredentials.UserName)

	// Отправляем ответ клиенту
	if _, err := io.WriteString(w, response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Устанавливаем успешный статус ответа
	w.WriteHeader(http.StatusOK)
}

// SaveUserNoteHandler обрабатывает запросы на сохранение заметки пользователя
func (h *handler) SaveUserNoteHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Создаем контекст для запроса
	ctx := context.Background()

	// Читаем тело запроса для получения данных заметки пользователя
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Распаковываем данные из тела запроса в структуру Note
	var requestNote models.Note
	if err := json.Unmarshal(buf.Bytes(), &requestNote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем обязательные поля (имя пользователя и заголовок заметки)
	if requestNote.UserName == "" || requestNote.Title == nil || *requestNote.Title == "" {
		http.Error(w, "user name and note title should not be empty", http.StatusBadRequest)
		return
	}

	// Сохраняем заметку пользователя в хранилище
	if err := h.db.SaveNote(ctx, requestNote); err != nil {
		message, status := handleUserError(requestNote.UserName, err)
		http.Error(w, message, status)
		return
	}

	// Формируем сообщение об успешном сохранении заметки
	response := fmt.Sprintf("Заметка для пользователя %q была успешно сохранена", requestNote.UserName)

	// Отправляем ответ клиенту
	if _, err := io.WriteString(w, response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Устанавливаем успешный статус ответа
	w.WriteHeader(http.StatusOK)
}

// GetUserNoteHandler обрабатывает запросы на получение заметки пользователя
func (h *handler) GetUserNoteHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Создаем контекст для запроса
	ctx := r.Context()

	// Читаем тело запроса для получения данных заметки пользователя
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Распаковываем данные из тела запроса в структуру Note
	var userNotesRequest models.Note
	if err := json.Unmarshal(body, &userNotesRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем заметку пользователя из хранилища goph-keeper
	creds, err := h.db.GetNotes(ctx, userNotesRequest)
	if err != nil {
		message, status := handleUserError(userNotesRequest.UserName, err)
		http.Error(w, message, status)
		return
	}

	// Преобразуем данные заметки пользователя в формат JSON
	notesResponse, err := json.Marshal(creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ клиенту
	if _, err = io.WriteString(w, string(notesResponse)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем статус успешного выполнения запроса
	w.WriteHeader(http.StatusOK)
}

// DeleteUserNotesHandler обрабатывает запросы на удаление заметок пользователя
func (h *handler) DeleteUserNotesHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Создаем контекст для запроса
	ctx := r.Context()

	// Читаем тело запроса для получения данных пользователя и заметки
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Распаковываем данные из тела запроса в структуру Note
	var userNotesRequest models.Note
	if err := json.Unmarshal(body, &userNotesRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Удаляем заметки пользователя из хранилища goph-keeper
	if err = h.db.DeleteNotes(ctx, userNotesRequest); err != nil {
		message, status := handleUserError(userNotesRequest.UserName, err)
		http.Error(w, message, status)
		return
	}

	// Составляем ответ

	response := fmt.Sprintf("Заметки для пользователя %q были успешно удалены", userNotesRequest.UserName)
	if userNotesRequest.Title != nil {
		response = fmt.Sprintf("Заметки для пользователя %q с заголовком %q были успешно удалены", userNotesRequest.UserName, *userNotesRequest.Title)
	}

	// Отправляем ответ клиенту
	if _, err := io.WriteString(w, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем статус успешного выполнения запроса
	w.WriteHeader(http.StatusOK)
}

// UpdateUserNoteHandler обрабатывает запросы на обновление заметки пользователя
func (h *handler) UpdateUserNoteHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Создаем контекст для запроса
	ctx := r.Context()

	// Читаем тело запроса
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Распаковываем данные из тела запроса в структуру Note
	var requestNote models.Note
	if err := json.Unmarshal(buf.Bytes(), &requestNote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, что заголовок и содержимое заметки не пустые
	if requestNote.Title == nil || requestNote.Content == nil {
		http.Error(w, "title and content should not be empty", http.StatusBadRequest)
		return
	}

	// Обновляем заметку пользователя в хранилище goph-keeper
	if err := h.db.UpdateNote(ctx, requestNote); err != nil {
		message, status := handleUserError(requestNote.UserName, err)
		http.Error(w, message, status)
		return
	}

	// Формируем ответ
	response := fmt.Sprintf("Заметка для пользователя %q успешно обновлена", requestNote.UserName)

	// Отправляем ответ клиенту
	if _, err := io.WriteString(w, response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Устанавливаем статус успешного выполнения запроса
	w.WriteHeader(http.StatusOK)
}

// SaveCardHandler обрабатывает запросы на сохранение карточки пользователя
func (h *handler) SaveCardHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

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

	// Устанавливаем статус успешного выполнения запроса
	w.WriteHeader(http.StatusOK)
}

// GetCardHandler обрабатывает запросы на получение карточек пользователя
func (h *handler) GetCardHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

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

	// Устанавливаем статус успешного выполнения запроса
	w.WriteHeader(http.StatusOK)
}

// DeleteCardHandler обрабатывает запросы на удаление карточек пользователя
func (h *handler) DeleteCardHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

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

	// Устанавливаем статус успешного выполнения запроса
	w.WriteHeader(http.StatusOK)
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
