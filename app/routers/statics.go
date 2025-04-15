package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/remote-serial-port-server/static"
)

func staticRoutesRegister(router *gin.RouterGroup) {
	router.StaticFS("/assets", static.GetFS(static.AssetsFS, "assets"))
	router.GET("/favicon.ico", func(context *gin.Context) {
		context.Writer.WriteHeader(200)
		icon, _ := static.RootFS.ReadFile("dist/icon.ico")
		_, _ = context.Writer.Write(icon)
		context.Writer.Header().Add("Accept", "image/x-icon")
	})

	router.GET("/", func(context *gin.Context) {
		context.Writer.WriteHeader(200)
		index, _ := static.RootFS.ReadFile("dist/index.html")
		_, _ = context.Writer.Write(index)
		context.Writer.Header().Add("Accept", "text/html")
	})
}
