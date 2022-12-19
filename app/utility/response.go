package utility

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func JsonResponse(code int, msg string, data gin.H, c *gin.Context) {
	if data != nil {
		c.JSON(code, gin.H{
			"msg":  msg,
			"data": data,
		})
	} else {
		c.JSON(code, gin.H{
			"msg": msg,
		})
	}
}

func JsonSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"code": 1,
		"msg":  "OK",
	})
}
