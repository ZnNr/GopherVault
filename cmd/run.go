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

// runCmd представляет команду run
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A command for running application server.",
	Run:   runHandler,
}

// runHandler обработчик команды запуска сервера
func runHandler(cmd *cobra.Command, args []string) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()
	defer logger.Sync() // Сброс буфера, если есть

	// Запуск приложения
	if err := run(logger.Sugar()); err != nil {
		log.Fatalf(err.Error())
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
}

// run инициализирует и запускает сервер приложения
func run(sugar *zap.SugaredLogger) error {
	var cfg models.Params
	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("Ошибка при загрузке переменных окружения: %w", err)
	}
	pg, err := database.New(cfg)
	if err != nil {
		return fmt.Errorf("Ошибка при попытке настройки БД: %w", err)
	}
	defer pg.Close()

	// Инициализация сервера
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.ApplicationPort))
	if err != nil {
		return fmt.Errorf("Ошибка при попытке прослушивания: %w", err)
	}
	router := router.New(pg, sugar)
	server := &http.Server{
		Handler: router,
	}
	go func() {
		if err := server.Serve(listener); err != nil {
			sugar.Errorf("Ошибка при запуске сервера: %v", err)
		}
	}()
	// Грациозное завершение работы
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-ch
		sugar.Infof(" Получен сигнал завершения работы сервера.")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			sugar.Infof("Не удалось корректно остановить сервер: %v", err)
		}
	}()

	// Перехват сигналов
	sugar.Infof("Started server on %s", cfg.ApplicationPort)
	ch = make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sugar.Infof(fmt.Sprint(<-ch))
	sugar.Infof("Stopping API server.")

	return nil
}
