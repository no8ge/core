package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/no8ge/core/pkg/handler"
)

func V1(r *gin.Engine) {

	core := r.Group("/v1")
	{
		core.GET("/", func(c *gin.Context) {
			log.Println("Hello, Geektutu")
			c.String(200, "Hello, Geektutu")
		})
		core.GET("/info", handler.Info)
	}

	webHook := r.Group("/v1")
	{
		webHook.POST("/validate", handler.Validate)
		webHook.POST("/inject", handler.Inject)
	}

}
