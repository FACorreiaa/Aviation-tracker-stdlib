package pprof

import (
	"context"
	"github.com/FACorreiaa/go-ollama/internal/logs"
	"github.com/gin-gonic/gin"
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

type Pprof interface {
	Run() error
	Shutdown(ctx context.Context) error
}

func New(config Config) Pprof {
	handler := gin.New()
	InitPprof(handler)
	return &server{
		port:      config.port,
		handler:   handler,
		certFile:  config.certFile,
		keyFile:   config.keyFile,
		enableTls: config.enableTls,
	}
}
