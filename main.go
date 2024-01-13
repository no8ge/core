package main

import (
	"github.com/gin-gonic/gin"
	"github.com/no8ge/core/api"
)

func main() {
	r := gin.Default()

	api.V1(r)

	// r.Run() // listen and serve on 0.0.0.0:8080
	r.RunTLS(":8080", "./certificates/core.pem", "./certificates/core.key")
}
