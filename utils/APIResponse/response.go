package APIResponse

import "github.com/gin-gonic/gin"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var C *gin.Context

func Error(message string) {
	if len(message) == 0 {
		message = "fail"
	}
	C.JSON(200, Response{
		Code:    -1,
		Message: message,
		Data:    nil,
	})
}
func Success(data interface{}) {
	C.JSON(200, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}
