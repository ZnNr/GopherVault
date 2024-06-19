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

// SaveUserNoteHandler обрабатывает запросы на сохранение заметки пользователя
func (h *handler) SaveUserNoteHandler(w http.ResponseWriter, r *http.Request) {
	h.cookiesMu.Lock()
	defer h.cookiesMu.Unlock()

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
}

// GetUserNoteHandler обрабатывает запросы на получение заметки пользователя
func (h *handler) GetUserNoteHandler(w http.ResponseWriter, r *http.Request) {
	h.cookiesMu.Lock()
	defer h.cookiesMu.Unlock()

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
}

// DeleteUserNotesHandler обрабатывает запросы на удаление заметок пользователя
func (h *handler) DeleteUserNotesHandler(w http.ResponseWriter, r *http.Request) {
	h.cookiesMu.Lock()
	defer h.cookiesMu.Unlock()

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

	response := fmt.Sprintf("Заметки для пользователя %q были успешно удалены", userNotesRequest.UserName)
	if userNotesRequest.Title != nil {
		response = fmt.Sprintf("Заметки для пользователя %q с заголовком %q были успешно удалены", userNotesRequest.UserName, *userNotesRequest.Title)
	}

	// Отправляем ответ клиенту
	if _, err := io.WriteString(w, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateUserNoteHandler обрабатывает запросы на обновление заметки пользователя
func (h *handler) UpdateUserNoteHandler(w http.ResponseWriter, r *http.Request) {
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
}
