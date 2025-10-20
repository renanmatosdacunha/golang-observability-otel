package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userRequest struct {
	Username string `json:"username" binding:"required"`
}

type userResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullname"`
}

var userOK = User{
	ID:       "1",
	Username: "Test",
	CPF:      "123456789101",
	FullName: "Test dos Testes",
}

func (server *Server) GetUser(ctx *gin.Context) {
	var req userRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logger.Error(
			"JSON Bind error",
			slog.String("error:", err.Error()),
		)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Username != userOK.Username {
		//server.logger.Error("user not found", slog.String("username:", req.Username))
		//server.logger.Error("user not found:", req.Username)
		logMessage := fmt.Sprintf("user not found, username: %s", req.Username)
		server.logger.Error(logMessage) // Não passa mais atributos aqui.
		ctx.JSON(http.StatusNotFound, "Error not found")
		return
	}

	user := userResponse{
		ID:       userOK.ID,
		Username: userOK.Username,
		FullName: userOK.FullName,
	}

	ctx.JSON(http.StatusOK, user)
}
