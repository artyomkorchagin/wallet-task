package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/artyomkorchagin/gin-boilerplate/config"
	"github.com/artyomkorchagin/gin-boilerplate/internal/logger"
	"github.com/artyomkorchagin/gin-boilerplate/internal/router"
	"github.com/spf13/viper"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func init() {
	config.LoadConfig()
}

//	@title			Comfortel Task
//	@version		1.0

//	@contact.name	Artyom Korchagin
//	@contact.email	artyomkorchagin333@gmail.com

//	@host		localhost:3000
//	@BasePath	/

func main() {
	var zapLogger *zap.Logger
	var err error

	if viper.GetString("ENV") == "development" {
		zapLogger, err = logger.NewDevelopmentLogger()
	} else {
		zapLogger, err = logger.NewLogger()
	}

	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("Starting application")

	db, err := sql.Open("pgx", config.GetDSN())
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := db.Ping(); err != nil {
		zapLogger.Fatal("Failed to ping database", zap.Error(err))
	}
	zapLogger.Info("Connected to database")

	if err := somepostgresql.RunMigrations(db); err != nil {
		zapLogger.Fatal("Failed to run up migration", zap.Error(err))
	}
	zapLogger.Info("Succesfully ran up migration")

	someRepo := somepostgresql.NewRepository(db)
	someSvc := someservice.NewService(someRepo)

	handler := router.NewHandler(userSvc, zapLogger)
	r := handler.InitRouter()

	srv := &http.Server{
		Addr:    ":" + viper.GetString("SERVER_PORT"),
		Handler: r,
	}

	go func() {
		zapLogger.Info("Server starting", zap.String("port", viper.GetString("SERVER_PORT")))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	zapLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error("Server shutdown failed", zap.Error(err))
	}

	zapLogger.Info("Server exited")

	if err := db.Close(); err != nil {
		zapLogger.Error("Error closing database connection", zap.Error(err))
	}
}
