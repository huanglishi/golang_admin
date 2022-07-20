package rmold

import (
	"basegin/app/common"

	"github.com/gin-gonic/gin"
)

//公共路由
func Apic(R *gin.Engine) {

	//表格
	auth := R.Group("/common/table")
	{
		auth.POST("/weigh", common.Weigh)
	}
	//上传图片
	uploadfile := R.Group("/common/uploadfile")
	{
		uploadfile.POST("/onefile", common.OneFile)
		uploadfile.GET("/getimage", common.GetImage)
		uploadfile.GET("/getimagebase", common.Getimagebase)
	}
	//测试
	testpath := R.Group("/common/api")
	{
		testpath.POST("/registry", common.Testpath)
	}
}
