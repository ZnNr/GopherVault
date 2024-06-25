package cmd

import (
	"fmt"
	cmdutil "github.com/ZnNr/GopherVault/cmdutils"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

// getNotesCmd представляет команду get-note
var getNotesCmd = &cobra.Command{
	Use:     "get-note",
	Short:   "Get user's notes from GopherVault",
	Example: "GopherVault get-note --user <user-name>",
	Run:     getNotesHandler,
}

func getNotesHandler(cmd *cobra.Command, args []string) {
	cfg := cmdutil.LoadEnvVariables()
	userName, title, _, _, _, _, _, _ := cmdutil.GetFlagsValues(cmd)

	requestNotes := models.Note{
		UserName: userName,
	}
	if title != "" {
		requestNotes.Title = &title
	}
	body := cmdutil.ConvertToJSONRequestNotes(requestNotes)

	resp, err := cmdutil.ExecutePostRequest(fmt.Sprintf("http://%s:%s/get/note", cfg.ApplicationHost, cfg.ApplicationPort), body)
	if err != nil {
		log.Printf(err.Error())
	}

	cmdutil.HandleResponse(resp, http.StatusOK)
}

func init() {
	rootCmd.AddCommand(getNotesCmd)
	getNotesCmd.Flags().String("user", "", "user name")
	getNotesCmd.Flags().String("title", "", "title of the note")
	getNotesCmd.MarkFlagRequired("user")
}
