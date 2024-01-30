package router

import (
	"github.com/gin-gonic/gin"
	"github.com/no8ge/core/internal/handler"
)

func V1(r *gin.Engine) {

	core := r.Group("/v1")
	{
		core.GET("/info", handler.Info)
	}

	webHook := r.Group("/v1")
	{
		webHook.POST("/validate", handler.Validate)
		webHook.POST("/inject", handler.Inject)
	}

}
