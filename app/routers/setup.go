package routers

import "github.com/gin-gonic/gin"

func SetupRouter(router *gin.RouterGroup) {
	debugRoutesRegister(router.Group("/debug"))
	apiRoutesRegister(router.Group("/api"))
	staticRoutesRegister(router)
}
