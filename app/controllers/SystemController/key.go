package SystemController

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/remote-serial-port-server/app/models"
	"github.com/iceking2nd/remote-serial-port-server/global"
	"net/http"
)

func GetAPIKey(c *gin.Context) {
	models.NewResponse(models.RESPONSE_OK, "ok", gin.H{"key": global.APIKey}).ResponseJson(http.StatusOK, c)
	return
}
