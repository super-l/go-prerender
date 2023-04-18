package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	resultData := make(map[string]interface{})
	resultData["version"] = "1.0.0"
	resultData["name"] = "go-prerender"

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    resultData,
		"message": "Welcome to use go-prerender!",
	})
}
