package models

import (
	"context"
)

// UserStorage интерфейс для работы с данными пользователей
//
//go:generate mockery --disable-version-string --filename user_storage_mock.go --name UserStorage
type UserStorage interface {
	SaveUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, userID string) (User, error)
	DeleteUser(ctx context.Context, userID string) error
	UpdateUser(ctx context.Context, user User) error
}

// NoteStorage интерфейс для работы с заметками
//
//go:generate mockery --disable-version-string --filename note_storage_mock.go --name NoteStorage
type NoteStorage interface {
	SaveNote(ctx context.Context, note Note) error
	GetNotes(ctx context.Context, noteRequest Note) ([]Note, error)
	DeleteNotes(ctx context.Context, noteRequest Note) error
	UpdateNote(ctx context.Context, note Note) error
}

// CardStorage интерфейс для работы с карточками
//
//go:generate mockery --disable-version-string --filename card_storage_mock.go --name CardStorage
type CardStorage interface {
	SaveCard(ctx context.Context, card Card) error
	GetCard(ctx context.Context, cardRequest Card) ([]Card, error)
	DeleteCards(ctx context.Context, cardRequest Card) error
	UpdateCard(ctx context.Context, cardUpdate Card) error
}

// CredentialsStorage интерфейс для работы с учетными данными
//
//go:generate mockery --disable-version-string --filename credentials_storage_mock.go --name CredentialsStorage
type CredentialsStorage interface {
	SaveCredentials(ctx context.Context, credentialsRequest Credentials) error
	GetCredentials(ctx context.Context, credentialsRequest Credentials) ([]Credentials, error)
	DeleteCredentials(ctx context.Context, credentialsRequest Credentials) error
	UpdateCredentials(ctx context.Context, credentials Credentials) error
}

// AuthenticationService определяет методы для регистрации пользователей, аутентификации и закрытия сервиса работы с аутентификацией.
//
//go:generate mockery --disable-version-string --filename authentication_service_mock.go --name AuthenticationService
type AuthenticationService interface {
	Register(ctx context.Context, login string, password string) error
	Login(ctx context.Context, login string, password string) error
	Close() error
}

//go:generate mockery --disable-version-string --filename storage_mock.go --name Storage
type Storage interface {
	SaveCredentials(ctx context.Context, credentialsRequest Credentials) error
	GetCredentials(ctx context.Context, credentialsRequest Credentials) ([]Credentials, error)
	DeleteCredentials(ctx context.Context, credentialsRequest Credentials) error
	UpdateCredentials(ctx context.Context, credentials Credentials) error
	SaveNote(ctx context.Context, note Note) error
	GetNotes(ctx context.Context, noteRequest Note) ([]Note, error)
	DeleteNotes(ctx context.Context, noteRequest Note) error
	UpdateNote(ctx context.Context, note Note) error
	SaveCard(ctx context.Context, card Card) error
	GetCard(ctx context.Context, cardRequest Card) ([]Card, error)
	DeleteCards(ctx context.Context, cardRequest Card) error
	Register(ctx context.Context, login string, password string) error
	Login(ctx context.Context, login string, password string) error
	Close() error
}
