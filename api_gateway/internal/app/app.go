package app

import (
	"api_gateway/config"
	grpc_client "api_gateway/internal/delivery/grpc/clients"
	api "api_gateway/internal/delivery/http"
	"api_gateway/pkg/http_server"
	"api_gateway/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	logger := logger.New(cfg.LoggerLevel)
	rpcConns, err := grpc_client.NewRPCClients(cfg)
	if err != nil {
		logger.Fatal("failed to connect to grpc server: %s", err.Error())
	}
	logger.Info("HTTP server started at: " + cfg.HTTPPort)
	handler := gin.New()
	api.NewRouter(handler, logger, rpcConns)

	httpServer := http_server.New(handler, http_server.Port(cfg.HTTPPort))

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-sig:
		logger.Info("app interrupt signal %s", s.String())
	case err := <-httpServer.Notify():
		logger.Error(" %s", err.Error())
	}
}
