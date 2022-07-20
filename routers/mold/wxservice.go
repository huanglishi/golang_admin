package rmold

import (
	"basegin/app/wxservice/signapi"
	"basegin/app/wxservice/wxconnect"

	"github.com/gin-gonic/gin"
)

//微信服务号
func ApiWxs(R *gin.Engine) {

	//微信管理
	wechat := R.Group("/wechat")
	{
		//服务端接口
		wxsetting_v1 := wechat.Group("/api/v1")
		{
			wxsetting_v1.POST("/push", wxconnect.WXMsgReceive)
			wxsetting_v1.GET("/push", wxconnect.WXCheckSignature)
		}
		//微信接口基础
		signapi_v1 := wechat.Group("/signapi/v1")
		{
			signapi_v1.GET("/jsdksignature", signapi.JSDKsignature)
		}

	}
}
