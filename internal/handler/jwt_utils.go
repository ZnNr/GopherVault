package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ZnNr/GopherVault/internal/database"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"strings"
	"time"
)

// parseUserInput парсирует входные данные пользователя и возвращает структуру *models.User
func parseUserInput(body io.ReadCloser) (*models.User, error) {
	var userFromRequest *models.User

	// Чтение данных из тела запроса
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(body); err != nil {
		return nil, fmt.Errorf("ошибка при чтении тела запроса: %w", err)
	}

	// Распаковка данных JSON в структуру *internal.User
	if err := json.Unmarshal(buf.Bytes(), &userFromRequest); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании JSON-данных: %w", err)
	}

	// Проверка наличия логина и пароля
	if userFromRequest.Login == "" || userFromRequest.Password == "" {
		return nil, fmt.Errorf("логин или пароль пустой")
	}

	return userFromRequest, nil
}

// извлечение JWT токена из строки куки
func extractJwtToken(cookie string) (*jwt.Token, error) {
	splitted := strings.Split(cookie, " ")
	if len(splitted) != 2 {
		return nil, ErrNoToken
	}

	tokenString := splitted[1]
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// создание JWT токена на основе имени пользователя и времени истечения
func createToken(username string, expiration time.Time) (string, error) {
	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// обработка ошибок, связанных с пользовательским запросом
func handleUserError(username string, err error) (string, int) {
	switch {
	case errors.Is(err, database.ErrNoSuchUser):
		return fmt.Sprintf("пользователя %q не существует", username), http.StatusUnauthorized
	case errors.Is(err, database.ErrInvalidCredentials):
		return fmt.Sprintf("предоставлен неверный пароль для пользователя %q", username), http.StatusUnauthorized
	case errors.Is(err, database.ErrUserAlreadyExists):
		return fmt.Sprintf("логин %q уже занят", username), http.StatusConflict
	case errors.Is(err, database.ErrNoData):
		return fmt.Sprintf("нет данных для пользователя %q", username), http.StatusNoContent
	case errors.Is(err, jwt.ErrSignatureInvalid), errors.Is(err, jwt.ErrTokenExpired), errors.Is(err, ErrTokenIsEmpty), errors.Is(err, ErrNoToken):
		return fmt.Sprintf("проблема с токеном для пользователя %q: %s", username, err.Error()), http.StatusUnauthorized
	default:
		return fmt.Sprintf("ошибка запроса пользователя %q: %s", username, err.Error()), http.StatusInternalServerError
	}
}
