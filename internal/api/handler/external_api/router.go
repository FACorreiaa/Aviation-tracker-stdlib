package external_api

import (
	"github.com/FACorreiaa/go-ollama/internal/api/handler/external_api/auth"
	"github.com/FACorreiaa/go-ollama/internal/api/handler/external_api/health"
	"github.com/FACorreiaa/go-ollama/internal/api/service"
	"github.com/gin-gonic/gin"
)

func initRouter(s *service.Service) *gin.Engine {
	router := gin.New()
	authHandler := auth.NewHandler(s)
	v1 := router.Group("/v1")
	{
		authGroup := v1.Group("/user")
		{
			authGroup.GET("/sign-up", authHandler.SignUp)
			authGroup.POST("/sign-in", authHandler.SignUp)
		}
	}

	health := new(health.HandlerHealth)

	router.GET("/health", health.Status)
	return router

}
