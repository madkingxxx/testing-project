package grpc_client

import (
	"api_gateway/config"
	fus "api_gateway/genproto/file_processing_service"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type rpcClients struct {
	fileProcessingServiceClient fus.FileProcessingServiceClient
}

func (r *rpcClients) FileProcessingClient() fus.FileProcessingServiceClient {
	return r.fileProcessingServiceClient
}

type RPCClients interface {
	FileProcessingClient() fus.FileProcessingServiceClient
}

func NewRPCClients(cfg *config.Config) (*rpcClients, error) {
	fileProcessingServiceConnection, err := grpc.Dial(cfg.FileServiceURL, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		return nil, errors.New("failed to connect to file processing service")
	}

	return &rpcClients{
		fileProcessingServiceClient: fus.NewFileProcessingServiceClient(fileProcessingServiceConnection),
	}, nil
}
