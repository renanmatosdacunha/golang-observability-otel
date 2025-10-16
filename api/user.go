package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type userRequest struct {
	Username string `json:"username" binding:"required"`
}

type userResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fulname"`
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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if req.Username != userOK.Username {
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
