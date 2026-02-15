package application

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sport-assistance/internal/handlers"
	"sport-assistance/internal/middlewares"
	"sport-assistance/internal/repositories"
	"sport-assistance/internal/services"
	"sport-assistance/pkg/configs"
	"sport-assistance/pkg/databases"
	"sport-assistance/pkg/logger"
	"sport-assistance/pkg/server"
	"syscall"
	"time"
)

type App struct {
	server *server.Server
	logger *slog.Logger
}

func NewApplication() *App {
	cfg := configs.GetConfigs()

	newLogger := logger.New(cfg.Logger)

	newLogger.Info("Connecting to database...")
	conn, err := databases.ConnectDB(cfg)
	if err != nil {
		newLogger.Error("Error connecting to database...")
	}

	newRepository := repositories.NewRepository(conn, newLogger)
	newRedisClient := databases.ConnectRedis(cfg)
	newService := services.NewService(newRepository, newLogger, cfg, newRedisClient)
	newMiddleware := middlewares.NewMiddleware(newRepository, cfg.SecurityConfig, newLogger, newRedisClient)
	newHandler := handlers.NewHandler(newService, newLogger, newMiddleware)
	newServer := server.NewServer(newHandler.InitHandler())

	return &App{
		server: newServer,
		logger: newLogger,
	}
}

func (a *App) Run() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		a.logger.Info("Starting application....")
		if err := a.server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Info(err.Error())
		}
	}()

	<-stopChan
	a.logger.Info("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("Server forced to shutdown: %v", err)
	}

	a.logger.Info("Server exited properly")
}
