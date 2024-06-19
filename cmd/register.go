package cmd

import (
	"fmt"
	cmdutil "github.com/ZnNr/GopherVault/cmdutils"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

// registerCmd представляет команду регистрации пользователей.
var registerCmd = &cobra.Command{
	Use:     "register",
	Short:   "Register in the GopherVault system.",
	Long:    `Register in the GopherVault system with provided login and password`,
	Example: "GopherVault register --login <user-system-login> --password <user-system-password>",
	Run:     registerHandler,
}

func registerHandler(cmd *cobra.Command, args []string) {
	cfg := cmdutil.LoadEnvVariables()
	login, password, _, _, _, _, _, _ := cmdutil.GetFlagsValues(cmd)

	userCreds := models.User{
		Login:    login,
		Password: password,
	}

	body := cmdutil.ConvertToJSONRequestUserCredential(userCreds)

	resp, err := cmdutil.ExecutePostRequest(fmt.Sprintf("http://%s:%s/auth/register", cfg.ApplicationHost, cfg.ApplicationPort), body)
	if err != nil {
		log.Printf(err.Error())
	}

	cmdutil.HandleResponse(resp, http.StatusOK)

	log.Printf("Пользователь %q успешно зарегистрировался в систему GopherVault\n", login)
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().String("login", "", "user login")
	registerCmd.Flags().String("password", "", "user password")
}
