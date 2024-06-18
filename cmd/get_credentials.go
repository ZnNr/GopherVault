package cmd

import (
	"fmt"
	cmdutil "github.com/ZnNr/GopherVault/cmdutils"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

// getCredentialsCmd представляет команду get-credentials
var getCredentialsCmd = &cobra.Command{
	Use:   "get-credentials",
	Short: "Get a pair of login/password for specified user",
	Long: `Get a pair of login/password for specified user from GopherVault storage. 
Only authorized users can use this command`,
	Example: "GopherVault get-credentials --user user_name",
	Run:     getCredentialsHandler,
}

func getCredentialsHandler(cmd *cobra.Command, args []string) {
	cfg := cmdutil.LoadEnvVariables()
	userName, userLogin, _, _, _, _, _ := cmdutil.GetFlagsValues(cmd)

	requestUserCredentials := models.Credentials{
		UserName: userName,
	}
	if userLogin != "" {
		requestUserCredentials.Login = &userLogin
	}
	body := cmdutil.ConvertToJSONRequestCredential(requestUserCredentials)

	resp, err := cmdutil.ExecutePostRequest(fmt.Sprintf("http://%s:%s/get/credentials", cfg.ApplicationHost, cfg.ApplicationPort), body)
	if err != nil {
		log.Printf(err.Error())
	}

	cmdutil.HandleResponse(resp, http.StatusOK)
}

func init() {
	rootCmd.AddCommand(getCredentialsCmd)
	getCredentialsCmd.Flags().String("user", "", "user name")
	getCredentialsCmd.Flags().String("login", "", "user login")
	getCredentialsCmd.MarkFlagRequired("user")
}
