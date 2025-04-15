package PortController

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/remote-serial-port-server/app/models"
	"go.bug.st/serial"
	"net/http"
)

func List(c *gin.Context) {
	//log := global.Log.WithField("function", "app.controllers.PortController.List")
	ports, err := serial.GetPortsList()
	if err != nil {
		models.NewResponse(models.RESPONSE_GET_PORTS_LIST_ERROR, err.Error(), nil)
		return
	}
	models.NewResponse(models.RESPONSE_OK, "ok", gin.H{"ports": ports}).ResponseJson(http.StatusOK, c)
	return
}
