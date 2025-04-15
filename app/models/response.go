package models

import "github.com/gin-gonic/gin"

const (
	RESPONSE_OK ResponseCode = iota
	RESPONSE_GET_PORTS_LIST_ERROR
	RESPONSE_BAD_BUADRATE
	RESPONSE_BAD_DATABITS
	RESPONSE_OPEN_SERIAL_PORT_ERROR
)

type ResponseCode int

type Response struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"msg"`
	Data    interface{}  `json:"data"`
}

func (r *Response) SetCode(code ResponseCode) *Response {
	r.Code = code
	return r
}

func (r *Response) SetMessage(msg string) *Response {
	r.Message = msg
	return r
}

func (r *Response) SetData(data interface{}) *Response {
	r.Data = data
	return r
}

func NewResponse(code ResponseCode, msg string, data interface{}) *Response {
	r := new(Response)
	return r.SetCode(code).SetMessage(msg).SetData(data)
}

func (r *Response) ResponseJson(httpCode int, c *gin.Context) {
	c.JSON(httpCode, r)
}
