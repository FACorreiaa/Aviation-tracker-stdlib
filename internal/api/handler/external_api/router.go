package external_api

import (
	"github.com/FACorreiaa/go-ollama/internal/api/handler/external_api/auth"
	"github.com/FACorreiaa/go-ollama/internal/api/handler/external_api/health"
	"github.com/FACorreiaa/go-ollama/internal/api/handler/external_api/index"
	"github.com/FACorreiaa/go-ollama/internal/api/service"
	"github.com/gin-gonic/gin"
	"html/template"
	"strings"
)

func initRouter(s *service.Service) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	authHandler := auth.NewHandler(s)
	router.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
	})
	router.Use(gin.Recovery())

	router.Static("/css", "./templates/css")
	router.LoadHTMLFiles("./templates/index.html")

	//router.LoadHTMLGlob("templates/*.html")
	v1 := router.Group("/v1")
	{
		authGroup := v1.Group("/user")
		{
			authGroup.GET("/sign-up", authHandler.SignUp)
			authGroup.POST("/sign-in", authHandler.SignUp)
		}

		indexGroup := v1.Group("/index")
		{
			indexGroup.GET("/", index.Home)
		}
	}

	health := new(health.HandlerHealth)

	router.GET("/health", health.Status)
	return router

}
