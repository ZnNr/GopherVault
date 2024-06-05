package cmd

import (
	"context"
	"fmt"
	"github.com/ZnNr/GopherVault/internal/database"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/ZnNr/GopherVault/internal/router"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A command for running application server.",
	Run: func(cmd *cobra.Command, args []string) {
		//init logger
		logger, _ := zap.NewProduction()
		defer logger.Sync() // flushes buffer, if any

		// run application
		if err := run(logger.Sugar()); err != nil {
			log.Fatalf(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

var credentialsStorage models.CredentialsStorage
var noteStorage models.NoteStorage
var cardStorage models.CardStorage
var authService models.AuthenticationService

func run(sugar *zap.SugaredLogger) error {
	var cfg models.Params
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("error while loading envs: %s\n", err)
	}
	pg, err := database.New(cfg)
	if err != nil {
		return fmt.Errorf("error while trying to setup DB: %w", err)
	}
	defer pg.Close()

	// init server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.ApplicationPort))
	if err != nil {
		return fmt.Errorf("error while trying to listen: %w", err)
	}
	router := router.New(credentialsStorage, noteStorage, cardStorage, authService, sugar)
	server := &http.Server{
		Handler: router,
	}
	go func() {
		server.Serve(listener)
	}()
	// graceful shutdown
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = server.Shutdown(ctx); err != nil {
			sugar.Infof("Could not shut down server correctly: %v\n", err)
			os.Exit(1)
		}
	}()

	// catch signals
	sugar.Infof("Started server on %s", cfg.ApplicationPort)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sugar.Infof(fmt.Sprint(<-ch))
	sugar.Infof("Stopping API server.")

	return nil
}
