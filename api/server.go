package api

import (
	"io"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	logger *slog.Logger
}

func NewServer(logger *slog.Logger) *Server {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	server := &Server{
		logger: logger,
	}
	router := gin.Default()

	router.GET("/users", server.GetUser)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {

	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
