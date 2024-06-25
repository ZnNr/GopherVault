package database

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ZnNr/GopherVault/internal/models"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"

	"strconv"
)

type Db struct {
	conn          *sql.DB
	encryptionKey string
	dataCipher    cipher.Block
}

// New создает новый экземпляр базы данных и возвращает его
func New(params models.Params) (*Db, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		params.StorageHost, params.StoragePort, params.StorageUser, params.StoragePassword, params.StorageDbName)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error while trying to open DB connection: %w", err)
	}
	c, err := aes.NewCipher([]byte(params.EncryptionKey))
	if err != nil {
		return nil, fmt.Errorf("error while creation cipher with key: %w", err)
	}
	pg := Db{
		conn:          conn,
		encryptionKey: params.EncryptionKey,
		dataCipher:    c,
	}

	if err = pg.conn.Ping(); err != nil {
		return nil, fmt.Errorf("error while trying to ping DB: %w", err)
	}
	return &pg, nil
}

// SaveNote сохраняет заметку в базе данных.
func (d *Db) SaveNote(ctx context.Context, noteRequest models.Note) error {
	encryptedContent, err := d.encryptAES(*noteRequest.Content)
	if err != nil {
		return fmt.Errorf("error encrypting your classified text: %w", err)
	}
	saveNotesQuery := "insert into notes (user_name, title, content, metadata) values ($1, $2, $3, $4)"
	if _, err = d.conn.ExecContext(ctx, saveNotesQuery, noteRequest.UserName, noteRequest.Title, encryptedContent, noteRequest.Metadata); err != nil {
		return fmt.Errorf("ошибка при сохранении заметки для пользователя %q: %w", noteRequest.UserName, err)
	}
	return nil
}

// GetNotes получает записи заметок из базы данных в соответствии с переданным запросом о заметках.
func (d *Db) GetNotes(ctx context.Context, noteRequest models.Note) ([]models.Note, error) {
	log.Printf("Выполняется запрос заметок для пользователя: %q", noteRequest.UserName)

	// Подготовка аргументов для запроса
	queryArgs := []interface{}{noteRequest.UserName}
	query := "select user_name, title, content, metadata from notes where user_name = $1"

	// Если указано название заметки, добавляем его в запрос и аргументы
	if noteRequest.Title != nil {
		queryArgs = append(queryArgs, *noteRequest.Title)
		query += fmt.Sprintf(" AND title = $%d", len(queryArgs))
	}

	// Инициируем запрос к базе данных
	rows, err := d.conn.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		errMsg := fmt.Errorf("ошибка при получении заметок для пользователя %q: %w", noteRequest.UserName, err)
		log.Println(errMsg)
		return nil, errMsg
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("ошибка при закрытии строк: %v", cerr)
		}
		if err := rows.Err(); err != nil {
			log.Printf("ошибка после итерации по результатам запроса: %v", err)
		}
	}()

	// Объявляем список для хранения заметок
	var notes []models.Note
	for rows.Next() {
		var userName, title, content string
		var metadata sql.NullString
		if err := rows.Scan(&userName, &title, &content, &metadata); err != nil {
			errMsg := fmt.Errorf("ошибка при сканировании строк после запроса на получение заметок пользователя: %w", err)
			log.Println(errMsg)
			return nil, errMsg
		}

		// Расшифровка контента заметки
		decryptedContent, err := d.decryptAES(content)
		if err != nil {
			errMsg := fmt.Errorf("ошибка при расшифровке контента заметки: %w", err)
			log.Println(errMsg)
			return nil, errMsg
		}

		// Создание структуры заметки и добавление в список
		note := models.Note{
			UserName: userName,
			Title:    &title,
			Content:  &decryptedContent,
		}
		if metadata.Valid {
			note.Metadata = &metadata.String
		}
		notes = append(notes, note)
	}

	if len(notes) == 0 {
		log.Printf("Для пользователя %q не найдено заметок", noteRequest.UserName)
		return nil, ErrNoData
	}
	log.Printf("Запрос заметок для пользователя %q выполнен успешно", noteRequest.UserName)

	return notes, nil
}

// DeleteNotes удаляет заметки из базы данных в соответствии с переданным запросом о заметке.
func (d *Db) DeleteNotes(ctx context.Context, noteRequest models.Note) error {
	// Подготовка аргументов для запроса
	args := []interface{}{noteRequest.UserName}
	deleteNotesQuery := "delete from notes where user_name = $1"

	// Добавляем критерий выборки по названию заметки, если он указан
	if noteRequest.Title != nil {
		args = append(args, *noteRequest.Title)
		deleteNotesQuery += " AND title = $2"
	}

	// Выполняем запрос на удаление заметок
	if _, err := d.conn.ExecContext(ctx, deleteNotesQuery, args...); err != nil {
		return fmt.Errorf("ошибка при удалении заметок для пользователя %q: %w", noteRequest.UserName, err)
	}
	return nil
}

// UpdateNote обновляет информацию о заметке в базе данных в соответствии с переданным запросом о заметке.
func (d *Db) UpdateNote(ctx context.Context, noteRequest models.Note) error {
	// Шифруем контент заметки
	encryptedContent, err := d.encryptAES(*noteRequest.Content)
	if err != nil {
		return fmt.Errorf("ошибка шифрования контента заметки: %w", err)
	}

	// Подготовка и выполнение запроса на обновление заметки
	updateNoteQuery := "update notes set content = $1, metadata = $2 where user_name = $3 and title = $4"
	if _, err := d.conn.ExecContext(ctx, updateNoteQuery, encryptedContent, noteRequest.Metadata, noteRequest.UserName, *noteRequest.Title); err != nil {
		return fmt.Errorf("ошибка при обновлении заметки %q для пользователя %q: %w", *noteRequest.Title, noteRequest.UserName, err)
	}
	return nil
}

// SaveCredentials сохраняет учетные данные в базе данных.
func (d *Db) SaveCredentials(ctx context.Context, credentialsRequest models.Credentials) error {
	// Шифруем пароль с использованием AES
	encryptedPassword, err := d.encryptAES(*credentialsRequest.Password)
	if err != nil {
		return fmt.Errorf("error encrypting your classified text: %w", err)
	}

	// Запрос для сохранения учетных данных
	saveCredsQuery := "insert into credentials (user_name, login, password, metadata) values ($1, $2, $3, $4)"
	_, err = d.conn.ExecContext(ctx, saveCredsQuery, credentialsRequest.UserName, *credentialsRequest.Login, encryptedPassword, credentialsRequest.Metadata)
	if err != nil {
		return fmt.Errorf("error while saving credentials for user %q: %w", credentialsRequest.UserName, err)
	}
	return nil
}

// GetCredentials получает учетные данные из базы данных.
func (d *Db) GetCredentials(ctx context.Context, credentialsRequest models.Credentials) ([]models.Credentials, error) {
	args := []interface{}{credentialsRequest.UserName}
	getCredsQuery := "select user_name, login, password, metadata from credentials where user_name = $1"
	if credentialsRequest.Login != nil {
		args = append(args, *credentialsRequest.Login)
		getCredsQuery += fmt.Sprintf(" AND login = $%d", len(args))
	}
	rows, err := d.conn.QueryContext(ctx, getCredsQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error while getting credentials for user %q: %w", credentialsRequest.UserName, err)
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	var creds []models.Credentials
	for rows.Next() {
		var userName, login, password string
		var metadata sql.NullString
		if err = rows.Scan(&userName, &login, &password, &metadata); err != nil {
			return nil, fmt.Errorf("error while scanning rows after get user credentials query: %w", err)
		}
		// Дешифруем пароль
		decryptedPassword, err := d.decryptAES(password)
		if err != nil {
			return nil, fmt.Errorf("error while decrypting password: %w", err)
		}
		res := models.Credentials{
			UserName: userName,
			Login:    &login,
			Password: &decryptedPassword,
		}
		if metadata.Valid {
			res.Metadata = &metadata.String
		}
		creds = append(creds, res)
	}
	if len(creds) == 0 {
		return nil, ErrNoData
	}
	return creds, nil
}

// DeleteCredentials удаляет учетные данные из базы данных.
func (d *Db) DeleteCredentials(ctx context.Context, credentialsRequest models.Credentials) error {
	args := []any{credentialsRequest.UserName}
	deleteCredsQuery := "delete from credentials where user_name = $1"
	if credentialsRequest.Login != nil {

		args = append(args, *credentialsRequest.Login)
		deleteCredsQuery += " AND login = $" + strconv.Itoa(len(args))
	}
	if _, err := d.conn.ExecContext(ctx, deleteCredsQuery, args...); err != nil {
		return fmt.Errorf("ошибка при удалении учетных данных для пользователя %q: %w", credentialsRequest.UserName, err)
	}
	return nil
}

// UpdateCredentials обновляет учетные данные в базе данных.
func (d *Db) UpdateCredentials(ctx context.Context, credentialsRequest models.Credentials) error {
	encryptedPassword, err := d.encryptAES(*credentialsRequest.Password)
	if err != nil {
		return fmt.Errorf("ошибка при шифровании пароля: %w", err)
	}
	updateCredsQuery := "update credentials set password = $1, metadata = $2 where user_name = $3 and login = $4"
	if _, err := d.conn.ExecContext(ctx, updateCredsQuery, encryptedPassword, credentialsRequest.Metadata, credentialsRequest.UserName, *credentialsRequest.Login); err != nil {
		return fmt.Errorf("ошибка при обновлении учетных данных для пользователя %q: %w", credentialsRequest.UserName, err)
	}
	return nil
}

// SaveCard сохраняет данные карты в базе данных.
func (d *Db) SaveCard(ctx context.Context, cardRequest models.Card) error {
	encryptedPassword, err := d.encryptAES(*cardRequest.Password)
	if err != nil {
		return fmt.Errorf("ошибка при шифровании пароля карты: %w", err)
	}
	encryptedCV, err := d.encryptAES(*cardRequest.CV)
	if err != nil {
		return fmt.Errorf("ошибка при шифровании CV карты: %w", err)
	}
	saveCardQuery := "insert into cards (user_name, bank_name, number, cv, password, cardType, metadata) values ($1, $2, $3, $4, $5, $6)"
	if _, err := d.conn.ExecContext(ctx, saveCardQuery, cardRequest.UserName, *cardRequest.BankName, *cardRequest.Number, encryptedCV, encryptedPassword, cardRequest.Metadata); err != nil {
		return fmt.Errorf("ошибка при сохранении данных карты для пользователя %q: %w", cardRequest.UserName, err)
	}
	return nil
}

// GetCard извлекает карты из базы данных на основе запроса.
func (d *Db) GetCard(ctx context.Context, cardRequest models.Card) ([]models.Card, error) {
	args := []interface{}{cardRequest.UserName}
	getCardsQuery := "select user_name, bank_name, number, cv, password, cardType, metadata from cards where user_name = $1"
	if cardRequest.BankName != nil {
		args = append(args, *cardRequest.BankName)
		getCardsQuery += fmt.Sprintf(" AND bank_name = $%d", len(args))
	}
	if cardRequest.Number != nil {
		args = append(args, *cardRequest.Number)
		getCardsQuery += fmt.Sprintf(" AND number = $%d", len(args))
	}
	rows, err := d.conn.QueryContext(ctx, getCardsQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении карт для пользователя %q: %w", cardRequest.UserName, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var cards []models.Card
	for rows.Next() {
		var userName, bankName, number, cv, password, cardType string
		var metadata sql.NullString
		if err = rows.Scan(&userName, &bankName, &number, &cv, &password, &cardType, &metadata); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строк после запроса на получение заметок пользователя: %w", err)
		}
		decryptedPassword, err := d.decryptAES(password)
		if err != nil {
			return nil, fmt.Errorf("ошибка при расшифровке пароля: %w", err)
		}
		decryptedCV, err := d.decryptAES(cv)
		if err != nil {
			return nil, fmt.Errorf("ошибка при расшифровке CV: %w", err)
		}
		res := models.Card{
			UserName: userName,
			BankName: &bankName,
			Number:   &number,
			CV:       &decryptedCV,
			Password: &decryptedPassword,
			CardType: &cardType,
		}
		if metadata.Valid {
			res.Metadata = &metadata.String
		}
		cards = append(cards, res)
	}
	if len(cards) == 0 {
		return nil, ErrNoData
	}
	return cards, nil
}

// DeleteCards удаляет карты из базы данных на основе запроса.
func (d *Db) DeleteCards(ctx context.Context, cardRequest models.Card) error {
	args := []interface{}{cardRequest.UserName}
	deleteNotesQuery := "delete from cards where user_name = $1"
	if cardRequest.Number != nil {
		args = append(args, *cardRequest.Number)
		deleteNotesQuery += fmt.Sprintf(" AND number = $%d", len(args))
	}
	if cardRequest.BankName != nil {
		args = append(args, *cardRequest.BankName)
		deleteNotesQuery += fmt.Sprintf(" AND bank_name = $%d", len(args))
	}
	if _, err := d.conn.ExecContext(ctx, deleteNotesQuery, args...); err != nil {
		return fmt.Errorf("ошибка при удалении карт для пользователя %q: %w", cardRequest.UserName, err)
	}
	return nil
}

// Login проверяет учетные данные пользователя в базе данных.
func (d *Db) Login(ctx context.Context, login string, password string) error {
	getRegisteredUser := `select login, password from registered_users where login = $1`

	var loginFromDB, passwordFromDB string
	if err := d.conn.QueryRowContext(ctx, getRegisteredUser, login).Scan(&loginFromDB, &passwordFromDB); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoSuchUser
		}
		return fmt.Errorf("ошибка при выполнении запроса на поиск: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordFromDB), []byte(password)); err != nil {
		return ErrInvalidCredentials
	}
	return nil
}

// Register добавляет нового пользователя в базу данных.
func (d *Db) Register(ctx context.Context, login string, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("этот пароль недопустим: %w", err)
	}
	registerUser := `insert into registered_users values ($1, $2)`
	if _, err = d.conn.ExecContext(ctx, registerUser, login, hash); err != nil {
		duplicateKeyErr := ErrDuplicateKey{Key: "registered_users_pkey"}
		if err.Error() == duplicateKeyErr.Error() {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("ошибка при выполнении запроса на регистрацию пользователя: %w", err)
	}
	return nil
}

// Close - метод для закрытия соединения с базой данных.
func (d *Db) Close() error {
	return d.conn.Close()
}

// encryptAES выполняет шифрование переданного текста с использованием AES и возвращает зашифрованный текст в виде base64 закодированной строки.
func (d *Db) encryptAES(plaintext string) (string, error) {
	// Инициализируем шифр с использованием режима CFB
	cfb := cipher.NewCFBEncrypter(d.dataCipher, []byte(d.encryptionKey)[:aes.BlockSize])

	// Выделяем память для зашифрованного текста
	cipherText := make([]byte, len(plaintext))

	// Шифруем переданный текст
	cfb.XORKeyStream(cipherText, []byte(plaintext))

	// Кодируем зашифрованный текст в base64 и возвращаем его
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// decryptAES расшифровывает переданный зашифрованный текст, используя AES, и возвращает исходный текст.
func (d *Db) decryptAES(ct string) (string, error) {
	// Декодируем base64 строку в байты зашифрованного текста
	cipherText, err := base64.StdEncoding.DecodeString(ct)
	if err != nil {
		return "", err
	}

	// Инициализируем дешифратор с использованием режима CFB
	cfb := cipher.NewCFBDecrypter(d.dataCipher, []byte(d.encryptionKey)[:aes.BlockSize])

	// Выделяем память для расшифрованного текста
	plainText := make([]byte, len(cipherText))

	// Расшифровываем текст
	cfb.XORKeyStream(plainText, cipherText)

	// Преобразуем байты расшифрованного текста в строку и возвращаем ее
	return string(plainText), nil
}
