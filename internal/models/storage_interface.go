package models

import (
	"context"
)

// Storage представляет интерфейс для работы с хранилищем
//
//go:generate mockery --disable-version-string --filename storage_mock.go --name Storage
type Storage interface {
	// SaveCredentials сохраняет учетные данные
	SaveCredentials(ctx context.Context, credentialsRequest Credentials) error

	// GetCredentials получает учетные данные
	GetCredentials(ctx context.Context, credentialsRequest Credentials) ([]Credentials, error)

	// DeleteCredentials удаляет учетные данные
	DeleteCredentials(ctx context.Context, credentialsRequest Credentials) error

	// UpdateCredentials обновляет учетные данные
	UpdateCredentials(ctx context.Context, credentials Credentials) error

	// SaveNote сохраняет заметку
	SaveNote(ctx context.Context, note Note) error

	// GetNotes получает заметки
	GetNotes(ctx context.Context, noteRequest Note) ([]Note, error)

	// DeleteNotes удаляет заметки
	DeleteNotes(ctx context.Context, noteRequest Note) error

	// UpdateNote обновляет заметку
	UpdateNote(ctx context.Context, note Note) error

	// SaveCard сохраняет карту
	SaveCard(ctx context.Context, card Card) error

	// GetCard получает карты
	GetCard(ctx context.Context, cardRequest Card) ([]Card, error)

	// DeleteCards удаляет карты
	DeleteCards(ctx context.Context, cardRequest Card) error

	// Register регистрирует пользователя
	Register(ctx context.Context, login string, password string) error

	// Login выполняет вход пользователя
	Login(ctx context.Context, login string, password string) error

	// Close закрывает соединение с хранилищем
	Close() error
}
