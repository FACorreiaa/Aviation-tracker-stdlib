package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/FACorreiaa/go-ollama/internal/api/handler"
	"github.com/FACorreiaa/go-ollama/internal/api/handler/external_api"
	"github.com/FACorreiaa/go-ollama/internal/api/handler/pprof"
	"github.com/FACorreiaa/go-ollama/internal/api/repository"
	"github.com/FACorreiaa/go-ollama/internal/api/repository/postgres"
	"github.com/FACorreiaa/go-ollama/internal/api/repository/redis"
	"github.com/FACorreiaa/go-ollama/internal/api/service"
	configs "github.com/FACorreiaa/go-ollama/internal/config"
	"github.com/FACorreiaa/go-ollama/internal/logs"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logs.InitDefaultLogger()

	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}

	config, err := configs.InitConfig(*environment)
	if err != nil {
		zap.L().Error("Config was not configure")
	}
	zap.L().Info("Config was successfully imported")

	zap.L().Info("Main logger was initialized successfully")

	if err := godotenv.Load(config.Dotenv); err != nil && config.Dotenv != "" {
		zap.L().Fatal("Dotenv was not loaded")
		os.Exit(1)
	}
	zap.L().Info("Dotenv file was successfully loaded")

	repositories := repository.NewRepository(
		repository.NewConfig(
			postgres.NewConfig(
				config.Repositories.Postgres.Scheme,
				config.Repositories.Postgres.Host,
				config.Repositories.Postgres.Port,
				config.Repositories.Postgres.Username,
				os.Getenv("DB_PASSWORD"),
				config.Repositories.Postgres.DB,
				config.Repositories.Postgres.SSLMode,
				time.Duration(config.Repositories.Postgres.MaxConnWaitingTime)*time.Second,
				postgres.CacheStatement,
			),
			redis.NewRedisConfig(
				os.Getenv("REDIS_HOST"),
				config.Repositories.Redis.RedisPassword,
				config.Repositories.Redis.RedisDb,
			),
		),
	)
	zap.L().Info("Repository was initialized")
	services := service.NewService(repositories)
	zap.L().Info("Service was initialized")
	handlers := handler.NewHandler(
		handler.NewConfig(
			external_api.NewConfig(
				config.Handlers.ExternalApi.Port,
				config.Handlers.ExternalApi.KeyFile,
				config.Handlers.ExternalApi.CertFile,
				config.Handlers.ExternalApi.EnableTLS,
			),
			pprof.NewConfig(
				config.Handlers.Pprof.Port,
				config.Handlers.Pprof.KeyFile,
				config.Handlers.Pprof.CertFile,
				config.Handlers.Pprof.EnableTLS,
			),
		),
		services,
	)
	zap.L().Info("Handler was initialized")

	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT,
	)
	var exitSignal os.Signal
	handlers.Handle(&exitSignal)
	zap.L().Info("Handler was successfully started")
	exitSignal = <-quit
	zap.L().Info("Exit...")
	handlers.Shutdown(context.Background())
	zap.L().Info("Handlers are shutdown")
}
