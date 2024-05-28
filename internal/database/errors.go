package database

import (
	"errors"
	"fmt"
)

// ErrDuplicateKey представляет ошибку дубликата ключа
type ErrDuplicateKey struct {
	Key string
}

// Error возвращает текст ошибки дубликата ключа
func (e ErrDuplicateKey) Error() string {
	return fmt.Sprintf("ERROR: duplicate key value violates unique constraint %q (SQLSTATE 23505)", e.Key)
}

// ErrUserAlreadyExists означает, что пользователь уже существует
var ErrUserAlreadyExists = errors.New("user already exists")

// ErrNoSuchUser означает, что пользователя не существует
var ErrNoSuchUser = errors.New("no such user")

// ErrInvalidCredentials означает неверные учетные данные
var ErrInvalidCredentials = errors.New("incorrect password")

// ErrNoData означает отсутствие данных для пользователя
var ErrNoData = errors.New("no data for user")
