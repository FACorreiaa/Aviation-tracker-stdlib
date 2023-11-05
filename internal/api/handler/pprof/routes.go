package pprof

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/pprof"
)

func InitPprof(router *gin.Engine) {
	prefixRouter := router.Group("/debug/pprof")
	{
		prefixRouter.GET("/", gin.WrapF(pprof.Index))
		prefixRouter.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		prefixRouter.GET("/profile", gin.WrapF(pprof.Profile))
		prefixRouter.POST("/symbol", gin.WrapF(pprof.Symbol))
		prefixRouter.GET("/symbol", gin.WrapF(pprof.Symbol))
		prefixRouter.GET("/trace", gin.WrapF(pprof.Trace))
		prefixRouter.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		prefixRouter.GET("/block", gin.WrapH(pprof.Handler("block")))
		prefixRouter.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		prefixRouter.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		prefixRouter.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		prefixRouter.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}
}

func ginWrapper(fn func(w http.ResponseWriter, r *http.Request)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fn(ctx.Writer, ctx.Request)
	}
}
