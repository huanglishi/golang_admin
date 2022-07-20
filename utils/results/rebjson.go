package results

import (
	"time"

	"github.com/gin-gonic/gin"
)

//请求成功的时候 使用该方法返回信息
func Success(ctx *gin.Context, msg string, data interface{}, token interface{}) {
	ctx.JSON(200, gin.H{
		"code":  0,
		"msg":   msg,
		"data":  data,
		"token": token,
		"time":  time.Now().Unix(),
	})
}

//请求失败的时候, 使用该方法返回信息
func Failed(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(200, gin.H{
		"code": 1,
		"msg":  msg,
		"data": data,
		"time": time.Now().Unix(),
	})

}
