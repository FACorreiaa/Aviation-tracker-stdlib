package handler

import (
	"context"
	"github.com/FACorreiaa/go-ollama/internal/api/handler/external_api"
	"github.com/FACorreiaa/go-ollama/internal/api/handler/pprof"
	"github.com/FACorreiaa/go-ollama/internal/api/service"
	"github.com/FACorreiaa/go-ollama/internal/logs"
	"os"
	"sync"
	"syscall"
)

type handler interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type Config struct {
	externalApiConfig external_api.Config
	pprofConfig       pprof.Config
}

func NewConfig(
	apiConfig external_api.Config,
	pprofConfig pprof.Config,
) Config {
	return Config{
		externalApiConfig: apiConfig,
		pprofConfig:       pprofConfig,
	}
}

type Handler struct {
	service *service.Service
	config  Config

	externalApi handler
	pprof       handler
}

func NewHandler(
	c Config,
	s *service.Service,
) *Handler {
	return &Handler{
		config:  c,
		service: s,
	}
}

func (h *Handler) Handle(exitSignal *os.Signal) {
	h.externalApi = external_api.New(h.config.externalApiConfig, h.service)
	h.pprof = pprof.New(h.config.pprofConfig)
	go func() {
		if err := h.pprof.Run(); err != nil && exitSignal == nil {
			logs.DefaultLogger.WithError(err).Fatal("Pprof server was closed unexpectedly")
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}
	}()
	go func() {
		if err := h.externalApi.Run(); err != nil && exitSignal == nil {
			logs.DefaultLogger.WithError(err).Fatal("REST API Server was closed unexpectedly")
			syscall.Kill(syscall.Getpid(), syscall.SIGQUIT)
		}
	}()
}

func (h *Handler) Shutdown(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		if err := h.externalApi.Shutdown(ctx); err != nil {
			logs.DefaultLogger.WithError(err).Fatal("Error on restApi shutdown")
		}
		wg.Done()
	}()
	go func() {
		if err := h.pprof.Shutdown(ctx); err != nil {
			logs.DefaultLogger.WithError(err).Fatal("Error on pprof shutdown")
		}
		wg.Done()
	}()
	wg.Wait()
}
