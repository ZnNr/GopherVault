package cmd

import (
	"encoding/json"
	"fmt"
	cmdutil "github.com/ZnNr/GopherVault/cmdutils"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/spf13/cobra"
	"log"
	"net/http"
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

func addNoteHandler(cmd *cobra.Command, args []string) {
	cfg := cmdutil.LoadEnvVariables()

	userName, title, content, metadata, _, _, _ := cmdutil.GetFlagsValues(cmd)

	requestNote := createNoteRequest(userName, title, content, metadata)

	body, err := json.Marshal(requestNote)
	if err != nil {
		log.Fatalf(err.Error())
	}

	resp, err := cmdutil.ExecutePostRequest(fmt.Sprintf("http://%s:%s/save/note", cfg.ApplicationHost, cfg.ApplicationPort), body)
	if err != nil {
		log.Printf(err.Error())
	}

	cmdutil.HandleResponse(resp, http.StatusOK)
}

func createNoteRequest(userName, title, content, metadata string) models.Note {
	requestNote := models.Note{
		UserName: userName,
		Title:    &title,
		Content:  &content,
	}
	if metadata != "" {
		requestNote.Metadata = &metadata
	}
	return requestNote
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
