package cmd

import (
	"fmt"
	cmdutil "github.com/ZnNr/GopherVault/cmdutils"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

// loginCmd представляет команду login
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the GopherVault system",
	Long: `Login to the GopherVault system with specified login and password. 
Only registered users can run this command`,
	Example: "GopherVault login --login <user-system-login> --password <user-system-password>`",
	Run:     loginHandler,
}

func loginHandler(cmd *cobra.Command, args []string) {
	cfg := cmdutil.LoadEnvVariables()
	login, password, _, _, _, _, _, _ := cmdutil.GetFlagsValues(cmd)

	userCreds := models.User{
		Login:    login,
		Password: password,
	}

	body := cmdutil.ConvertToJSONRequestUserCredential(userCreds)

	resp, err := cmdutil.ExecutePostRequest(fmt.Sprintf("http://%s:%s/auth/login", cfg.ApplicationHost, cfg.ApplicationPort), body)
	if err != nil {
		log.Printf(err.Error())
	}

	cmdutil.HandleResponse(resp, http.StatusOK)

	log.Printf("Пользователь %q успешно вошел в систему GopherVault\n", login)
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().String("login", "", "user login")
	loginCmd.Flags().String("password", "", "user password")
	loginCmd.MarkFlagRequired("login")
	loginCmd.MarkFlagRequired("password")
}
