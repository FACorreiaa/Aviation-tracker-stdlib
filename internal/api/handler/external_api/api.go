package external_api

import (
	"context"
	"github.com/FACorreiaa/go-ollama/internal/api/service"
	"github.com/FACorreiaa/go-ollama/internal/logs"
	"syscall"
)

type Config struct {
	port      string
	keyFile   string
	certFile  string
	enableTls bool
}

func NewConfig(
	port string,
	keyFile string,
	certFile string,
	enableTls bool,
) Config {
	if enableTls && (certFile == "" || keyFile == "") {
		logs.DefaultLogger.Fatal("Tls is enabled but cert file or key file doesn't have a path")
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}
	return Config{
		port:      port,
		keyFile:   keyFile,
		certFile:  certFile,
		enableTls: enableTls,
	}
}

type Api interface {
	Run() error
	Shutdown(ctx context.Context) error
}

func New(config Config, s *service.Service) Api {
	return &server{
		port:      config.port,
		handler:   initRouter(s),
		certFile:  config.certFile,
		keyFile:   config.keyFile,
		enableTls: config.enableTls,
	}
}
