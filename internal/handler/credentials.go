package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ZnNr/GopherVault/internal/models"
	"io"
	"net/http"
)

// GetUserCredentialsHandler обрабатывает запросы на получение учетных данных пользователя
func (h *handler) GetUserCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	h.cookiesMu.Lock()
	defer h.cookiesMu.Unlock()

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
}

// SaveUserCredentialsHandler обрабатывает запросы на сохранение учетных данных пользователя
func (h *handler) SaveUserCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	h.cookiesMu.Lock()
	defer h.cookiesMu.Unlock()

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
}

// DeleteUserCredentialsHandler обрабатывает запросы на удаление учетных данных пользователя
func (h *handler) DeleteUserCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	h.cookiesMu.Lock()
	defer h.cookiesMu.Unlock()

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
}

// UpdateUserCredentialsHandler обрабатывает запросы на обновление учетных данных пользователя
func (h *handler) UpdateUserCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	h.cookiesMu.Lock()
	defer h.cookiesMu.Unlock()

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
}
