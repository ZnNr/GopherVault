package cmd

import (
	"fmt"
	cmdutil "github.com/ZnNr/GopherVault/cmdutils"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

// addCredentialsCmd представляет команду add-credentials
var addCredentialsCmd = &cobra.Command{
	Use:   "add-credentials",
	Short: "Add a pair of login/password to GopherVault.",
	Long: `Add a pair of login/password to GopherVault database for
long-term storage. Only authorized users can use this command. The password is stored in the database in encrypted form.`,
	Example: "GopherVault add-credentials --user <user-name> --login <user-login> --password <password to store> --metadata <some description>",
	Run:     addCredentialsHandler,
}

func addCredentialsHandler(cmd *cobra.Command, args []string) {
	cfg := cmdutil.LoadEnvVariables()

	// Получение значений флагов из командной строки
	userName, login, password, metadata, _, _, _ := cmdutil.GetFlagsValues(cmd)

	requestCredentials := models.Credentials{
		UserName: userName,
		Login:    &login,
		Password: &password,
	}
	if metadata != "" {
		requestCredentials.Metadata = &metadata
	}

	body := cmdutil.ConvertToJSONRequestCredential(requestCredentials)

	resp, err := cmdutil.ExecutePostRequest(fmt.Sprintf("http://%s:%s/save/credentials", cfg.ApplicationHost, cfg.ApplicationPort), body)
	if err != nil {
		log.Printf(err.Error())
	}

	cmdutil.HandleResponse(resp, http.StatusOK)
}

func init() {
	rootCmd.AddCommand(addCredentialsCmd)
	addCredentialsCmd.Flags().String("user", "", "user name")
	addCredentialsCmd.Flags().String("login", "", "user login")
	addCredentialsCmd.Flags().String("password", "", "user password")
	addCredentialsCmd.Flags().String("metadata", "", "metadata")
	addCredentialsCmd.MarkFlagRequired("user")
	addCredentialsCmd.MarkFlagRequired("login")
	addCredentialsCmd.MarkFlagRequired("password")
}
