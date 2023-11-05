package auth

import (
	"github.com/FACorreiaa/go-ollama/internal/api/service"
	"github.com/FACorreiaa/go-ollama/internal/api/structs"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) SignUp(ctx *gin.Context) {
	_, _ = h.service.User.Create(structs.User{Username: "test"})
}
