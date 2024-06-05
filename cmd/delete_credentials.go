package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

// deleteCredentialsCmd represents the deleteCredentials command
var deleteCredentialsCmd = &cobra.Command{
	Use:     "delete-credentials",
	Short:   "Delete credentials for user from GopherVault storage",
	Example: "GopherVault delete-credentials --user <user-name> --login <user-login>",
	Run: func(cmd *cobra.Command, args []string) {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("error while getting envs: %s", err)
		}
		var cfg models.Params
		if err := envconfig.Process("", &cfg); err != nil {
			log.Fatalf("error while loading envs: %s\n", err)
		}

		userName, _ := cmd.Flags().GetString("user")
		login, _ := cmd.Flags().GetString("login")
		requestUserCredentials := models.Credentials{
			UserName: userName,
		}
		if login != "" {
			requestUserCredentials.Login = &login
		}
		body, err := json.Marshal(requestUserCredentials)
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Println(string(body))

		resp, err := resty.New().R().
			SetHeader("Content-type", "application/json").
			SetBody(body).
			Post(fmt.Sprintf("http://%s:%s/delete/credentials", cfg.ApplicationHost, cfg.ApplicationPort))
		if err != nil {
			log.Printf(err.Error())
		}
		if resp.StatusCode() != http.StatusOK {
			log.Printf("status code is not OK: %s\n", resp.Status())
		}
		log.Printf(resp.String())
	},
}

func init() {
	rootCmd.AddCommand(deleteCredentialsCmd)
	deleteCredentialsCmd.Flags().String("user", "", "user name")
	deleteCredentialsCmd.Flags().String("login", "", "user login")
	deleteCredentialsCmd.MarkFlagRequired("user")
}
