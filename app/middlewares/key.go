package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/remote-serial-port-server/global"
	"net/http"
)

func CheckAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-Key") != global.APIKey && c.DefaultQuery("key", "") != global.APIKey {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
