package PortController

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/iceking2nd/remote-serial-port-server/app/models"
	"github.com/iceking2nd/remote-serial-port-server/global"
	"github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Open(c *gin.Context) {
	log := global.Log.WithField("function", "app.controllers.PortController.Open")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket:", err)
		return
	}
	defer conn.Close()

	port := c.DefaultQuery("port", "")
	baudrate, err := strconv.Atoi(c.DefaultQuery("baudrate", "9600"))
	if err != nil {
		models.NewResponse(models.RESPONSE_BAD_BUADRATE, err.Error(), nil).ResponseJson(http.StatusOK, c)
		return
	}
	databits, err := strconv.Atoi(c.DefaultQuery("databits", "8"))
	if err != nil {
		models.NewResponse(models.RESPONSE_BAD_DATABITS, err.Error(), nil).ResponseJson(http.StatusOK, c)
		return
	}
	var parity serial.Parity
	switch c.DefaultQuery("parity", "none") {
	case "odd":
		parity = serial.OddParity
	case "even":
		parity = serial.EvenParity
	case "mark":
		parity = serial.MarkParity
	case "space":
		parity = serial.SpaceParity
	default:
		parity = serial.NoParity
	}
	var stopbits serial.StopBits
	switch c.DefaultQuery("stopbits", "1") {
	case "1.5":
		stopbits = serial.OnePointFiveStopBits
	case "2":
		stopbits = serial.TwoStopBits
	default:
		stopbits = serial.OneStopBit
	}
	rts := c.DefaultQuery("rts", "0") == "0"
	dtr := c.DefaultQuery("dtr", "0") == "0"

	log.WithFields(logrus.Fields{
		"port":     port,
		"baudrate": baudrate,
		"databits": databits,
		"parity":   parity,
		"stopbits": stopbits,
		"rts":      rts,
		"dtr":      dtr,
	}).Debug("open port with parameters")

	mode := &serial.Mode{
		BaudRate:          baudrate,
		DataBits:          databits,
		Parity:            parity,
		StopBits:          stopbits,
		InitialStatusBits: &serial.ModemOutputBits{rts, dtr},
	}

	serialPort, err := serial.Open(port, mode)
	if err != nil {
		models.NewResponse(models.RESPONSE_OPEN_SERIAL_PORT_ERROR, err.Error(), nil).ResponseJson(http.StatusOK, c)
		return
	}
	defer serialPort.Close()

	// 启动两个协程：一个用于读取串口数据，另一个用于接收 WebSocket 消息
	go func() {
		buf := make([]byte, 128)
		for {
			n, err := serialPort.Read(buf)
			if err != nil {
				log.WithError(err).Errorln("error reading from serial port")
				return
			}
			conn.WriteMessage(websocket.TextMessage, buf[:n])
		}
	}()

	// 监听 WebSocket 消息并发送到串口
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.WithError(err).Errorln("WebSocket read error")
			break
		}
		_, err = serialPort.Write(msg)
		if err != nil {
			log.WithError(err).Errorln("error writing to serial port")
			break
		}
	}
}
