package grpc

import (
	"bytes"
	"context"
	fus "file_processing_service/genproto/file_processing_service"
	"file_processing_service/internal/entity"
	"file_processing_service/internal/usecase/file_processing"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	fuc file_processing.FileUsecase
	fus.UnimplementedFileProcessingServiceServer
}

func (s *server) UploadFileWithStreaming(stream fus.FileProcessingService_UploadFileWithStreamingServer) error {
	chunk, err := stream.Recv()
	if err != nil {
		return err
	}
	file := &entity.File{
		Name:      chunk.GetName(),
		Content:   bytes.NewBuffer(nil),
		CreatedAt: time.Now().UTC(),
	}
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		_, err = file.Content.Write(chunk.GetData())
		if err != nil {
			return err
		}
	}
	err = s.fuc.UploadFile(file)
	if err != nil {
		return err
	}
	return stream.SendAndClose(&fus.FileUploadResponse{
		Message: "file uploaded successfully",
	})
}

func (s *server) DownloadFile(in *fus.DownloadFileRequest, stream fus.FileProcessingService_DownloadFileServer) error {
	file, err := s.fuc.DownloadFile(in.GetName())
	if err != nil {
		return err
	}
	const chunkSize = 2 << 10
	buffer := make([]byte, chunkSize)
	for {
		n, err := file.Content.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = stream.Send(&fus.Chunk{
			Data: buffer[:n],
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *server) GetFileList(ctx context.Context, in *emptypb.Empty) (*fus.FileList, error) {
	files, err := s.fuc.GetFileList()
	if err != nil {
		return nil, err
	}
	fileList := &fus.FileList{}
	for _, file := range files {
		fileList.FileInfo = append(fileList.FileInfo, &fus.FileInfo{
			Name:      file.Name,
			CreatedAt: file.CreatedAt.String(),
			UpdatedAt: file.UpdatedAt.String(),
		})
	}
	return fileList, nil
}

func NewFileUploadService(gservice *grpc.Server, fuc file_processing.FileUsecase) {
	fus.RegisterFileProcessingServiceServer(gservice, &server{
		fuc: fuc,
	})
	reflection.Register(gservice)
}
