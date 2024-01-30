package handler

import "github.com/gin-gonic/gin"

func Info(c *gin.Context) {
	resp := gin.H{
		"version": "1.0.0",
	}
	c.JSON(200, resp)

}
