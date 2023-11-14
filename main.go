package main

import (
	"context"
	"fmt"
	"github.com/FACorreiaa/go-ollama/config"
	"github.com/FACorreiaa/go-ollama/controller"
	"github.com/FACorreiaa/go-ollama/db"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var logHandler slog.Handler
	logHandlerOptions := slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.Log.Level,
	}
	if cfg.Log.Format == "json" {
		logHandler = slog.NewJSONHandler(os.Stdout, &logHandlerOptions)
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &logHandlerOptions)
	}
	slog.SetDefault(slog.New(logHandler))

	pool, err := db.Init(cfg.Database.ConnectionURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer pool.Close()

	db.WaitForDB(pool)

	if err = db.Migrate(pool); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		WriteTimeout: cfg.Server.WriteTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		Handler:      controller.Router(pool, []byte(cfg.Server.SessionKey)),
	}

	go func() {
		slog.Info("Starting server " + cfg.Server.Addr)
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("ListenAndServe", "error", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer cancel()
	srv.Shutdown(ctx)
	slog.Info("shutting down")
	os.Exit(0)
}
