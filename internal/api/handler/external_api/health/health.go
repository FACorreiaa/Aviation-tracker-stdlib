package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerHealth struct{}

func (h HandlerHealth) Status(c *gin.Context) {
	c.String(http.StatusOK, "Working!")
}
