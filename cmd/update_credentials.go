package cmd

import (
	"fmt"
	cmdutil "github.com/ZnNr/GopherVault/cmdutils"
	"github.com/ZnNr/GopherVault/internal/models"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// updateCredentialsCmd представляет команду updateCredentials
var updateCredentialsCmd = &cobra.Command{
	Use:     "update-credentials",
	Short:   "Update user credentials for the provided login.",
	Example: "GopherVault update-credentials --user <user-name> --login <saved-login> --password <new-password>",
	Run:     updateCredentialsHandler,
}

// updateCredentialsHandler обработчик команды обновления учетных данных
func updateCredentialsHandler(cmd *cobra.Command, args []string) {
	cfg := cmdutil.LoadEnvVariables()
	userName, login, password, metadata, _, _, _ := cmdutil.GetFlagsValues(cmd)

	requestCredentials := models.Credentials{
		UserName: userName,
		Login:    &login,
		Password: &password,
		Metadata: &metadata,
	}

	body := cmdutil.ConvertToJSONRequestCredential(requestCredentials)

	resp, err := cmdutil.ExecutePostRequest(fmt.Sprintf("http://%s:%s/update/credentials", cfg.ApplicationHost, cfg.ApplicationPort), body)
	if err != nil {
		log.Printf(err.Error())
	}

	cmdutil.HandleResponse(resp, http.StatusOK)
}

func init() {
	rootCmd.AddCommand(updateCredentialsCmd)
	updateCredentialsCmd.Flags().String("user", "", "user name")
	updateCredentialsCmd.Flags().String("login", "", "user login")
	updateCredentialsCmd.Flags().String("password", "", "user password")
	updateCredentialsCmd.Flags().String("metadata", "", "metadata")
	updateCredentialsCmd.MarkFlagRequired("user")
	updateCredentialsCmd.MarkFlagRequired("login")
	updateCredentialsCmd.MarkFlagRequired("password")
}
