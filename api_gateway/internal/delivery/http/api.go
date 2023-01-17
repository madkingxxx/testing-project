package api

import (
	grpc_client "api_gateway/internal/delivery/grpc/clients"
	"api_gateway/internal/delivery/http/handlers"
	"api_gateway/internal/delivery/http/middleware"
	"api_gateway/pkg/logger"

	_ "api_gateway/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8000
// @BasePath    /v1
// @schemes http https
func NewRouter(handler *gin.Engine, l logger.Interface, rpcConns grpc_client.RPCClients) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	handler.Use(middleware.LogMiddleware(l))
	h := handler.Group("/v1")
	{
		handlers.NewFileUploadRoutes(h, l, rpcConns)
	}
}
