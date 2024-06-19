package cmd

import (
	"fmt"
	cmdutil "github.com/ZnNr/GopherVault/cmdutils"
	"github.com/ZnNr/GopherVault/internal/models"
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
	cfg := cmdutil.LoadEnvVariables()
	userName, title, _, _, _, _, _, _ := cmdutil.GetFlagsValues(cmd)

	requestNotes := models.Note{
		UserName: userName,
	}
	if title != "" {
		requestNotes.Title = &title
	}
	body := cmdutil.ConvertToJSONRequestNotes(requestNotes)

	resp, err := cmdutil.ExecutePostRequest(fmt.Sprintf("http://%s:%s/delete/note", cfg.ApplicationHost, cfg.ApplicationPort), body)
	if err != nil {
		log.Printf(err.Error())
	}

	cmdutil.HandleResponse(resp, http.StatusOK)
}

func init() {
	rootCmd.AddCommand(deleteNotesCmd)
	deleteNotesCmd.Flags().String("user", "", "user name")
	deleteNotesCmd.Flags().String("title", "", "title of the note")
	deleteNotesCmd.MarkFlagRequired("user")
}
