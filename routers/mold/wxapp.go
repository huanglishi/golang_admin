package rmold

import (
	"basegin/app/wxapp/articleapi"
	"basegin/app/wxapp/docappointment"
	"basegin/app/wxapp/wxapi"
	"basegin/app/wxapp/wxcommon"
	"basegin/app/wxapp/wxhome"
	"basegin/app/wxapp/wxmall"
	"basegin/app/wxapp/wxorder"

	"github.com/gin-gonic/gin"
)

func Apiwxapp(R *gin.Engine) {
	home := R.Group("/wxapp")
	{
		//微信公账号接口
		wxapi_v1 := home.Group("/wxapi/v1")
		{
			wxapi_v1.GET("/getwechataccount", wxapi.Getwechataccount_v1)
			wxapi_v1.GET("/getaccesstoken", wxapi.Getaccesstoken_v1)
			wxapi_v1.POST("/refreshtoken", wxapi.Refreshtoken_v1)
			wxapi_v1.POST("/commonlogin", wxapi.Commonlogin_v1)
		}
		//公共方法
		common_v1 := home.Group("/common/v1")
		{
			common_v1.GET("/getarticledetail", wxcommon.Getarticledetail_v1)
			common_v1.POST("/dotag", wxcommon.Dotag_v1)    //关注，收藏
			common_v1.GET("/getfile", wxcommon.Getfile_v1) //获取附件
		}
		//首页
		home_v1 := home.Group("/home/v1")
		{
			home_v1.GET("/getbasetop", wxhome.Getbasetop_v1)
		}
		//文章
		article_v1 := home.Group("/article/v1")
		{
			article_v1.GET("/getarticlecate", articleapi.Getarticlecate_v1)
			article_v1.GET("/getarticlelist", articleapi.Getarticlelist_v1)
			article_v1.GET("/getarticledetail", articleapi.Getarticledetail_v1)
		}
		//医生预约
		docappointment_v1 := home.Group("/docappointment/v1")
		{
			docappointment_v1.GET("/getbase", docappointment.Getbase_v1)
			docappointment_v1.GET("/getlist", docappointment.Getlist_v1)
			docappointment_v1.POST("/getdetail", docappointment.Getdetail_v1)
			docappointment_v1.GET("/getmydoctor", docappointment.Getmydoctor_v1)
		}
		//商城
		wxmall_v1 := home.Group("/mall/v1")
		{
			wxmall_v1.GET("/getbase", wxmall.Getbase_v1)
			wxmall_v1.GET("/getlist", wxmall.Getlist_v1)
			wxmall_v1.GET("/getdetail", wxmall.Getdetail_v1)
			wxmall_v1.POST("/addcart", wxmall.Addcart_v1)
			wxmall_v1.GET("/getcart", wxmall.Getcart_v1)
			wxmall_v1.DELETE("/delcart", wxmall.Delcart_v1)
		}
		//订单
		order_v1 := home.Group("/order/v1")
		{
			order_v1.POST("/addorder", wxorder.Addorder_v1)
		}
		// //聊天
		// wechat_v1 := home.Group("/wechat/v1")
		// {
		// 	wechat_v1.POST("/getemojiDB", wechat.GetemojiDB_v1)
		// }
	}
}
