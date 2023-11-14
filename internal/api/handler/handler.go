package handler

import (
	"context"
	"github.com/FACorreiaa/go-ollama/internal/api/handler/external_api"
	"github.com/FACorreiaa/go-ollama/internal/api/handler/pprof"
	"github.com/FACorreiaa/go-ollama/internal/api/service"
	"go.uber.org/zap"
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
	//r := gin.New()
	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	//r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	//r.Use(gin.Recovery())

	h.externalApi = external_api.New(h.config.externalApiConfig, h.service)
	h.pprof = pprof.New(h.config.pprofConfig)
	go func() {
		if err := h.pprof.Run(); err != nil && exitSignal == nil {
			zap.L().Fatal("Pprof server was closed unexpectedly")
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}
	}()
	go func() {
		if err := h.externalApi.Run(); err != nil && exitSignal == nil {
			zap.L().Fatal("REST API Server was closed unexpectedly")
			syscall.Kill(syscall.Getpid(), syscall.SIGQUIT)
		}
	}()
}

func (h *Handler) Shutdown(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		if err := h.externalApi.Shutdown(ctx); err != nil {
			zap.L().Fatal("Error on restApi shutdown")
		}
		wg.Done()
	}()
	go func() {
		if err := h.pprof.Shutdown(ctx); err != nil {
			zap.L().Fatal("Error on pprof shutdown")
		}
		wg.Done()
	}()
	wg.Wait()
}
