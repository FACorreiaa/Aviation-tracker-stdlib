package pprof

import (
	"context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
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
		zap.L().Fatal("Tls is enabled but cert file or key file doesn't have a path")
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
	handler := mux.NewRouter()
	InitPprof(handler)
	return &server{
		port:      config.port,
		handler:   handler,
		certFile:  config.certFile,
		keyFile:   config.keyFile,
		enableTls: config.enableTls,
	}
}
