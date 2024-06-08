package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// addNotesCmd represents the add-notes command
var addNotesCmd = &cobra.Command{
	Use:   "add-note",
	Short: "Add user's note to GopherVault storage.",
	Long: `Add user's note to GopherVault database for long-term storage.
Only authorized users can use this command. The note content is stored in the database in encrypted form.`,
	Example: "GopherVault add-note --user <user-name> --title <note title> --content <note content> --metadata <note metadata>",
	Run:     addNoteHandler,
}

// addNoteHandler обработчик команды добавления заметки
func addNoteHandler(cmd *cobra.Command, args []string) {
	// Загружаем переменные среды из файла .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error while getting envs: %s", err)
	}

	var cfg models.Params
	// Загружаем переменные среды в структуру cfg
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Error while loading envs: %s\n", err)
	}

	// Получаем значения флагов команды
	userName, _ := cmd.Flags().GetString("user")
	title, _ := cmd.Flags().GetString("title")
	content, _ := cmd.Flags().GetString("content")
	metadata, _ := cmd.Flags().GetString("metadata")

	// Создаем объект модели Note
	requestNote := models.Note{
		UserName: userName,
		Title:    &title,
		Content:  &content,
	}
	// Добавляем метаданные, если они указаны
	if metadata != "" {
		requestNote.Metadata = &metadata
	}

	// Преобразуем объект Note в JSON
	body, err := json.Marshal(requestNote)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Отправляем POST-запрос на сохранение заметки
	resp, err := resty.New().R().
		SetHeader("Content-type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("http://%s:%s/save/note", cfg.ApplicationHost, cfg.ApplicationPort))
	if err != nil {
		log.Printf(err.Error())
	}

	// Проверяем статус ответа
	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code is not OK: %s\n", resp.Status())
	}
	fmt.Println(resp.String())
}

func init() {
	// Добавляем команду addNotesCmd к корневой команде
	rootCmd.AddCommand(addNotesCmd)
	// Добавляем флаги для команды addNotesCmd
	addNotesCmd.Flags().String("user", "", "user name")
	addNotesCmd.Flags().String("title", "", "user login")
	addNotesCmd.Flags().String("content", "", "user password")
	addNotesCmd.Flags().String("metadata", "", "metadata")
	// Помечаем флаги, как обязательные
	addNotesCmd.MarkFlagRequired("user")
	addNotesCmd.MarkFlagRequired("title")
	addNotesCmd.MarkFlagRequired("content")
}
