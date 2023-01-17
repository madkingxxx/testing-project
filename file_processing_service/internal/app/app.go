package app

import (
	"file_processing_service/config"
	fgrpc "file_processing_service/internal/delivery/grpc"
	"file_processing_service/internal/entity"
	"file_processing_service/internal/usecase/file_processing"
	"file_processing_service/pkg/logger"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func Run(cfg *config.Config) {
	logger := logger.New(cfg.LoggerLevel)

	r := entity.NewFileStorage()

	fileUC := file_processing.NewService(r)

	l, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		logger.Fatal(err)
	}

	gserver := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			fgrpc.GetUnaryInterceptors(logger, cfg.MaxUnaryRequestCount)...,
		),
		grpc.ChainStreamInterceptor(
			fgrpc.GetStreamingInterceptors(logger)...,
		),
		grpc.MaxConcurrentStreams(uint32(cfg.MaxStreamRequestCount)),
	)
	fgrpc.NewFileUploadService(gserver, fileUC)
	go func() {
		logger.Info("gRPC server started at: " + cfg.GRPCPort)
		if err := gserver.Serve(l); err != nil {
			logger.Fatal(err)
		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	s := <-sig
	logger.Info("received signal: ", s)
	gserver.GracefulStop()
}
