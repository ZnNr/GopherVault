package database

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"fmt"
	"github.com/ZnNr/GopherVault/internal/models"
)

type db struct {
	conn          *sql.DB
	encryptionKey string
	dataCipher    cipher.Block
}

// New создает новый экземпляр базы данных и возвращает его
func New(params models.Params) (*db, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		params.StorageHost, params.StoragePort, params.StorageUser, params.StoragePassword, params.StorageDbName)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error while trying to open DB connection: %w", err)
	}

	encryptionKeyBytes := []byte(params.EncryptionKey)
	dataCipher, err := aes.NewCipher(encryptionKeyBytes)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error while creating cipher with key: %w", err)
	}

	pg := &db{
		conn:          conn,
		encryptionKey: params.EncryptionKey,
		dataCipher:    dataCipher,
	}

	if err = pg.conn.Ping(); err != nil {
		pg.conn.Close()
		return nil, fmt.Errorf("error while trying to ping DB: %w", err)
	}

	return pg, nil
}
