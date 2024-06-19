package cmd

import (
	"fmt"
	cmdutil "github.com/ZnNr/GopherVault/cmdutils"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

// deleteCredentialsCmd представляет команду deleteCredentials
var deleteCredentialsCmd = &cobra.Command{
	Use:     "delete-credentials",
	Short:   "Delete credentials for user from GopherVault storage",
	Example: "GopherVault delete-credentials --user <user-name> --login <user-login>",
	Run:     deleteCredentialsHandler,
}

func deleteCredentialsHandler(cmd *cobra.Command, args []string) {

	cfg := cmdutil.LoadEnvVariables()
	userName, login, _, _, _, _, _, _ := cmdutil.GetFlagsValues(cmd)

	// Создаем объект модели Credentials для запроса
	requestUserCredentials := models.Credentials{
		UserName: userName,
	}
	// Добавляем логин, если он указан
	if login != "" {
		requestUserCredentials.Login = &login
	}

	body := cmdutil.ConvertToJSONRequestCredential(requestUserCredentials)

	// Отправляем POST-запрос на удаление учетных данных
	resp, err := cmdutil.ExecutePostRequest(fmt.Sprintf("http://%s:%s/delete/credentials", cfg.ApplicationHost, cfg.ApplicationPort), body)
	if err != nil {
		log.Printf(err.Error())
	}

	cmdutil.HandleResponse(resp, http.StatusOK)
}

func init() {
	rootCmd.AddCommand(deleteCredentialsCmd)
	deleteCredentialsCmd.Flags().String("user", "", "user name")
	deleteCredentialsCmd.Flags().String("login", "", "user login")
	deleteCredentialsCmd.MarkFlagRequired("user")
}
