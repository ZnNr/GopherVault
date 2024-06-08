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

// deleteNotesCmd представляет команду delete-note.
var deleteNotesCmd = &cobra.Command{
	Use:     "delete-note",
	Short:   "Delete user's notes from GopherVault storage",
	Example: "GopherVault delete-note --user <user-name> --title <note title>",
	Run:     deleteNotesHandler,
}

func deleteNotesHandler(cmd *cobra.Command, args []string) {
	// Загрузка переменных окружения
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Ошибка при загрузке окружения: %s", err)
	}

	var cfg models.Params
	// Загрузка настроек из переменных окружения
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Ошибка при загрузке окружения: %s\n", err)
	}

	userName, _ := cmd.Flags().GetString("user")
	title, _ := cmd.Flags().GetString("title")
	requestNotes := models.Note{
		UserName: userName,
	}
	if title != "" {
		requestNotes.Title = &title
	}
	body, err := json.Marshal(requestNotes)
	if err != nil {
		log.Fatalln(err.Error())
	}

	resp, err := resty.New().R().
		SetHeader("Content-type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("http://%s:%s/delete/note", cfg.ApplicationHost, cfg.ApplicationPort))
	if err != nil {
		log.Printf(err.Error())
	}
	if resp.StatusCode() != http.StatusOK {
		log.Printf("Статус ответа не ОК: %s\n", resp.Status())
	}
	log.Printf(resp.String())

}

func init() {
	rootCmd.AddCommand(deleteNotesCmd)
	deleteNotesCmd.Flags().String("user", "", "user name")
	deleteNotesCmd.Flags().String("title", "", "title of the note")
	deleteNotesCmd.MarkFlagRequired("user")
}
