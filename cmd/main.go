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

	"github.com/artyomkorchagin/wallet-task/config"
	"github.com/artyomkorchagin/wallet-task/internal/infrastructure"
	"github.com/artyomkorchagin/wallet-task/internal/logger"
	walletpostgresql "github.com/artyomkorchagin/wallet-task/internal/repository/postgres/wallet"
	"github.com/artyomkorchagin/wallet-task/internal/router"
	walletservice "github.com/artyomkorchagin/wallet-task/internal/services/wallet"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//	@title			Comfortel Task
//	@version		1.0

//	@contact.name	Artyom Korchagin
//	@contact.email	artyomkorchagin333@gmail.com

//	@host		localhost:3000
//	@BasePath	/

func main() {
	var err error

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	var zapLogger *zap.Logger

	if cfg.LogLevel == "DEV" {
		zapLogger, err = logger.NewDevelopmentLogger()
	} else {
		zapLogger, err = logger.NewLogger()
	}

	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("Starting application")
	zapLogger.Info("Connecting to database")
	db, err := sql.Open("pgx", config.GetDSN())
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}

	zapLogger.Info("Pinging database")
	if err := db.Ping(); err != nil {
		zapLogger.Fatal("Failed to ping database", zap.Error(err))
	}
	zapLogger.Info("Connected to database")

	zapLogger.Info("Running migrations")
	if err := walletpostgresql.RunMigrations(db); err != nil {
		zapLogger.Fatal("Failed to run up migration", zap.Error(err))
	}
	zapLogger.Info("Succesfully ran up migration")

	zapLogger.Info("Connecting to redis")
	rAddr := cfg.Redis.Host + cfg.Redis.Port
	rdb, err := infrastructure.NewRedisClient(rAddr, cfg.Redis.Password)
	if err != nil {
		zapLogger.Fatal("Failed to connect to redis", zap.Error(err))
	}

	zapLogger.Info("Connected to redis")

	walletRepo := walletpostgresql.NewRepository(db)
	walletSvc := walletservice.NewService(walletRepo, rdb, zapLogger)

	handler := router.NewHandler(walletSvc, zapLogger)
	r := handler.InitRouter()

	port := cfg.Server.Port
	srv := &http.Server{
		Addr:    cfg.Server.Host + ":" + port,
		Handler: r,
	}

	go func() {
		zapLogger.Info("Server starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	zapLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error("Server shutdown failed", zap.Error(err))
	}

	zapLogger.Info("Server exited")

	if err := db.Close(); err != nil {
		zapLogger.Error("Error closing database connection", zap.Error(err))
	}

	zapLogger.Info("Shutdown completed")
}
