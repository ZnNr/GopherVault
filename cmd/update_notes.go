package cmd

import (
	"fmt"
	cmdutil "github.com/ZnNr/GopherVault/cmdutils"
	"github.com/ZnNr/GopherVault/internal/models"
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
	cfg := cmdutil.LoadEnvVariables()
	userName, title, content, metadata, _, _, _ := cmdutil.GetFlagsValues(cmd)

	// Создаем объект заметки
	requestNote := models.Note{
		UserName: userName,
		Title:    &title,
		Content:  &content,
		Metadata: &metadata,
	}

	body := cmdutil.ConvertToJSONRequestNotes(requestNote)

	// Отправляем POST-запрос на сервер
	resp, err := cmdutil.ExecutePostRequest(fmt.Sprintf("http://%s:%s/update/note", cfg.ApplicationHost, cfg.ApplicationPort), body)
	if err != nil {
		log.Printf(err.Error())
	}

	cmdutil.HandleResponse(resp, http.StatusOK)
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
