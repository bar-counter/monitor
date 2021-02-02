package pprof

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

func OnBind(r *gin.RouterGroup, pprofPrefix string, handlers ...gin.HandlerFunc) {
	var pprofGroup *gin.RouterGroup
	if len(handlers) > 0 {
		handlerFunc := handlers[0]
		if handlerFunc != nil {
			pprofGroup = r.Group(pprofPrefix, handlerFunc)
		} else {
			pprofGroup = r.Group(pprofPrefix)
		}
	} else {
		pprofGroup = r.Group(pprofPrefix)
	}

	{
		pprofGroup.GET("/", pprofBizHandler(pprof.Index))
		pprofGroup.GET("/cmdline", pprofBizHandler(pprof.Cmdline))
		pprofGroup.GET("/profile", pprofBizHandler(pprof.Profile))
		pprofGroup.POST("/symbol", pprofBizHandler(pprof.Symbol))
		pprofGroup.GET("/symbol", pprofBizHandler(pprof.Symbol))
		pprofGroup.GET("/trace", pprofBizHandler(pprof.Trace))
		pprofGroup.GET("/allocs", pprofBizHandler(pprof.Handler("allocs").ServeHTTP))
		pprofGroup.GET("/block", pprofBizHandler(pprof.Handler("block").ServeHTTP))
		pprofGroup.GET("/goroutine", pprofBizHandler(pprof.Handler("goroutine").ServeHTTP))
		pprofGroup.GET("/heap", pprofBizHandler(pprof.Handler("heap").ServeHTTP))
		pprofGroup.GET("/mutex", pprofBizHandler(pprof.Handler("mutex").ServeHTTP))
		pprofGroup.GET("/threadcreate", pprofBizHandler(pprof.Handler("threadcreate").ServeHTTP))
	}
}

func pprofBizHandler(h http.HandlerFunc) gin.HandlerFunc {
	handler := http.HandlerFunc(h)
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
