package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

// updateNotesCmd представляет команду updateNotes
var updateNotesCmd = &cobra.Command{
	Use:     "update-note",
	Short:   "Update user notes.",
	Example: "GopherVault update-notes --user <user-name> --title <note-title> --content <new-content>",
	Run:     updateNoteHandler,
}

// updateNoteHandler обработчик команды обновления заметки
func updateNoteHandler(cmd *cobra.Command, args []string) {
	// Загружаем переменные окружения из файла .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Ошибка при загрузке переменных окружения: %s", err)
	}

	// Получаем конфигурацию
	var cfg models.Params
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %s\n", err)
	}

	// Получаем значения флагов из командной строки
	userName, _ := cmd.Flags().GetString("user")
	title, _ := cmd.Flags().GetString("title")
	content, _ := cmd.Flags().GetString("content")
	metadata, _ := cmd.Flags().GetString("metadata")

	// Создаем объект заметки
	requestNote := models.Note{
		UserName: userName,
		Title:    &title,
		Content:  &content,
		Metadata: &metadata,
	}

	// Преобразуем объект заметки в JSON
	body, err := json.Marshal(requestNote)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Отправляем POST-запрос на сервер
	resp, err := sendUpdateRequest(cfg, body)
	if err != nil {
		log.Printf(err.Error())
	}

	// Проверяем статус ответа
	if resp.StatusCode() != http.StatusOK {
		log.Printf("Статус ответа не 'OK': %s\n", resp.Status())
	}

	log.Println(resp.String())
}

// sendUpdateRequest отправляет запрос на обновление заметки
func sendUpdateRequest(cfg models.Params, body []byte) (*resty.Response, error) {
	return resty.New().R().
		SetHeader("Content-type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("http://%s:%s/update/note", cfg.ApplicationHost, cfg.ApplicationPort))
}

func init() {
	rootCmd.AddCommand(updateNotesCmd)
	// Добавляем флаги команды updateNotesCmd
	updateNotesCmd.Flags().String("user", "", "user name")
	updateNotesCmd.Flags().String("title", "", "title of the note")
	updateNotesCmd.Flags().String("content", "", "new note's content")
	updateNotesCmd.Flags().String("metadata", "", "metadata")
	updateNotesCmd.MarkFlagRequired("user")
	updateNotesCmd.MarkFlagRequired("title")
	updateNotesCmd.MarkFlagRequired("content")
}
