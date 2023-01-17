package handlers

import (
	"api_gateway/genproto/file_processing_service"
	grpc_client "api_gateway/internal/delivery/grpc/clients"
	"api_gateway/pkg/logger"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

type fileUploadRoutes struct {
	l          logger.Interface
	rpcClients grpc_client.RPCClients
}

// UploadFile godoc
// @Summary      Upload file
// @Description  Uploading large/small files
// @Tags         file
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "file"
// @Success      200  {object}  string
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /file/upload [post]
func (f *fileUploadRoutes) uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	reader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	client, err := f.rpcClients.FileProcessingClient().UploadFileWithStreaming(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer client.CloseSend()

	// initial request to send the file info
	err = client.Send(&file_processing_service.Chunk{
		Name: file.Filename,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	const chunkSize = 2 << 10
	buffer := make([]byte, chunkSize)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			f.l.Error("error reading chunk to buffer: %s", err.Error())
		}
		req := &file_processing_service.Chunk{
			Data: buffer[:n],
		}
		err = client.Send(req)
		if err != nil {
			f.l.Error("error sending chunk to server: %s", err.Error())
		}
	}

	res, err := client.CloseAndRecv()
	if err != nil && err != io.EOF {
		f.l.Error("error closing and receiving: %s", err.Error())
		c.JSON(http.StatusInternalServerError, fmt.Errorf("close and recieve error %w", err).Error())
		return
	}
	f.l.Info("file uploaded: %s", res.GetMessage())
	c.JSON(http.StatusOK, res.GetMessage())
}

// File list godoc
// @Summary      List of files
// @Description  Get list of files
// @Tags         file
// @Accept       application/json
// @Produce      json
// @Success      200  {object}  string
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /file/list [get]
func (f *fileUploadRoutes) listFile(c *gin.Context) {
	files, err := f.rpcClients.FileProcessingClient().GetFileList(c.Request.Context(), &emptypb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, files)
}

// Download file godoc
// @Summary      Download a file by name
// @Description  Download a file by name
// @Tags         file
// @Accept       application/json
// @Produce      json
// @Param        name path string true "file name"
// @Success      200  {file}    multipart/form-data
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /file/{name} [get]
func (f *fileUploadRoutes) downloadFile(c *gin.Context) {
	fileName := c.Param("name")
	client, err := f.rpcClients.FileProcessingClient().DownloadFile(c.Request.Context(), &file_processing_service.DownloadFileRequest{Name: fileName})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer client.CloseSend()
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/octet-stream")
	for {
		chunk, err := client.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			f.l.Error("error receiving chunk from server: %s", err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		_, err = c.Writer.Write(chunk.GetData())
		if err != nil {
			f.l.Error("error writing chunk to response: %s", err.Error())
		}
	}
	c.Status(http.StatusOK)
}

func NewFileUploadRoutes(handler *gin.RouterGroup, l logger.Interface, rpcConns grpc_client.RPCClients) {
	fr := &fileUploadRoutes{l, rpcConns}
	h := handler.Group("/file")
	{
		h.GET("/:name", fr.downloadFile)
		h.POST("/upload", fr.uploadFile)
		h.GET("/list", fr.listFile)
	}
}
