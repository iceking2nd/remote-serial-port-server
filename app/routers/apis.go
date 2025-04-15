package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/remote-serial-port-server/app/controllers/PortController"
	"github.com/iceking2nd/remote-serial-port-server/app/controllers/SystemController"
	"github.com/iceking2nd/remote-serial-port-server/app/middlewares"
)

func apiRoutesRegister(route *gin.RouterGroup) {
	systemRoutes := route.Group("/system")
	systemRoutes.GET("/key", SystemController.GetAPIKey)

	portRoutes := route.Group("/port")
	portRoutes.GET("/", middlewares.CheckAPIKey(), PortController.List)
	portRoutes.GET("/open", middlewares.CheckAPIKey(), PortController.Open)
}
