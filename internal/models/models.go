package models

import "github.com/golang-jwt/jwt/v4"

type UserBase struct {
	UserName string `json:"user_name"` // Имя пользователя
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Credentials struct {
	UserBase
	Login    *string `json:"login,omitempty"`    // Логин пользователя
	Password *string `json:"password,omitempty"` // Пароль пользователя
	Metadata *string `json:"metadata,omitempty"` // Дополнительная метаинформация
}

type Note struct {
	UserBase
	Title    *string `json:"title,omitempty"`    // Заголовок заметки
	Content  *string `json:"content,omitempty"`  // Содержание заметки
	Metadata *string `json:"metadata,omitempty"` // Дополнительная метаинформация
}

type Card struct {
	UserBase
	BankName *string `json:"bank_name,omitempty"` // Наименование банка
	Number   *string `json:"number,omitempty"`    // Номер карты
	CV       *string `json:"cv,omitempty"`        // Код CV (Security code)
	Password *string `json:"password,omitempty"`  // Пароль карты
	Metadata *string `json:"metadata,omitempty"`  // Дополнительная метаинформация
}

type Params struct {
	StoragePort     string `envconfig:"POSTGRES_PORT"`
	StorageHost     string `envconfig:"POSTGRES_HOST"`
	StorageUser     string `envconfig:"POSTGRES_USER"`
	StoragePassword string `envconfig:"POSTGRES_PASSWORD"`
	StorageDbName   string `envconfig:"POSTGRES_DB"`
	ApplicationPort string `envconfig:"APPLICATION_PORT"`
	ApplicationHost string `envconfig:"APPLICATION_HOST"`
	EncryptionKey   string `envconfig:"KEEPER_ENCRYPTION_KEY"`
}
