package rmold

import (
	"basegin/app/admin/auths"
	"basegin/app/admin/business"
	"basegin/app/admin/developers"
	"basegin/app/admin/user"
	"basegin/app/cmdshell"

	"github.com/gin-gonic/gin"
)

func Api(R *gin.Engine) {

	// 定义路由组
	// user := R.Group("/api/v1")
	// //
	// user.Use(middleware.Privilege())
	// {

	// }
	//登录路由
	adminr := R.Group("/admin/user")
	{ // 请求参数在请求路径上
		// id := c.Query("id") //查询请求URL后面拼接的参数
		// name := c.PostForm("name") //从表单中查询参数
		// uuid := c.Param("uuid") //取得URL中参数
		adminr.POST("/add", user.AddParam)
		adminr.GET("/info", user.GetInfo)
		adminr.POST("/login", user.Lonin)
		adminr.POST("/refreshtoken", user.Refreshtoken)
		adminr.POST("/logout", user.Logout)
		adminr.POST("/updata", user.Updata)
		adminr.GET("/list", user.QueryParam)
		adminr.DELETE("/del", user.DelParam)
		adminr.POST("/changepassword", user.Changepassword)
	}
	//权限路由
	auth := R.Group("/auth/rule")
	{
		auth.GET("/getlist", auths.Getlist)
		auth.GET("/getparent", auths.Getparent)
		auth.POST("/add", auths.Add)
		auth.POST("/lock", auths.Add_lock)
		auth.DELETE("/del", auths.Del)
	}
	//分组
	group := R.Group("/auth/group")
	{
		group.GET("/getlist", auths.Getlist_group)
		group.GET("/getparent", auths.Getparent_group)
		group.GET("/grouptree", auths.Getgroup_tree)
		group.POST("/add", auths.Add_group)
		group.DELETE("/del", auths.Del_group)
	}
	//业务管理-开户管理
	business_r := R.Group("/business")
	{
		//模块管理
		openaccount := business_r.Group("/openaccount")
		{
			openaccount.POST("/add", business.Add_cate)
			openaccount.GET("/getparentcate", business.Getparent_cate)
			openaccount.GET("/getcate", business.Get_cate)
			openaccount.DELETE("/delcate", business.Del_cate)
			openaccount.POST("/upcate", business.Up_cate)
			openaccount.GET("/getlist", business.Getlist)
			openaccount.POST("/addcontent", business.Add_content)
			openaccount.DELETE("/delcontent", business.Del_content)
			openaccount.POST("/upcontent", business.Up_content)
		}
		//套餐
		funpackage := business_r.Group("/funpackage")
		{
			funpackage.GET("/getmarketdata", business.Getmarketdata)
			funpackage.GET("/versionview", business.Versionview)
			funpackage.POST("/savecar", business.Savecar)
			funpackage.GET("/choosepackage", business.Choosepackage)
			funpackage.DELETE("/delchoose", business.Delchoose)
			funpackage.POST("/upchoose", business.Upchoose)
		}
		//网站
		website := business_r.Group("/website")
		{
			website.GET("/list", business.Websitelist)
			website.POST("/upaudit", business.Up_audit)
			website.GET("/getsitefield", business.Getsitefield_v1)
			website.POST("/upsite", business.Upsite_v1)
			website.GET("/sitelist", business.Sitelist_v1)
			website.DELETE("/clearsite", business.Clearsite_v1)
			website.POST("/clonesite", business.Clonesite_v1)
		}
	}

	//开发者
	developer := R.Group("/developer")
	{
		//模块管理
		module := developer.Group("/module")
		{
			module.POST("/add", developers.Add_cate)
			module.GET("/getcate", developers.Get_cate)
			module.GET("/getparent", developers.Getparent)
			module.DELETE("/delcate", developers.Del_cate)
			module.POST("/upcate", developers.Up_cate)
			module.GET("/getlist", developers.Getlist)
			module.POST("/addcontent", developers.Add_content)
			module.DELETE("/delcontent", developers.Del_content)
			module.POST("/upcontent", developers.Up_content)
			//版本
			module.POST("/addversion", developers.Addversion)
			module.GET("/getversion", developers.Getversion)
			module.DELETE("/delversion", developers.Delversion)
			module.POST("/upversion", developers.Upversion)
			module.GET("/getversiondes", developers.Getversiondes)
			module.POST("/upversiondes", developers.Upversiondes)
			module.GET("/getversionmenu", developers.Getversionmenu)
			module.POST("/upversionmenu", developers.Upversionmenu)
			module.GET("/getvisiturl", developers.Getvisiturl)
			module.POST("/upvisiturl", developers.Upvisiturl)
		}
		//接口文档
		apidoc := developer.Group("/apidoc")
		{
			apidoc.POST("/add", developers.Add_cate_apidoc)
			apidoc.GET("/getparentcate", developers.Getparent_cate)
			apidoc.GET("/getcate", developers.Get_cate_apidoc)
			apidoc.DELETE("/delcate", developers.Del_cate_apidoc)
			apidoc.POST("/upcate", developers.Up_cate_apidoc)
			apidoc.GET("/getlist", developers.Getlist_apidoc)
			apidoc.POST("/addcontent", developers.Add_content_apidoc)
			apidoc.DELETE("/delcontent", developers.Del_content_apidoc)
			apidoc.POST("/upcontent", developers.Up_content_apidoc)
			apidoc.GET("/getviewapidoc", developers.Getviewapidoc)
		}
		//后台菜单
		clientmenu := developer.Group("/clientmenu")
		{
			clientmenu.GET("/getlist", developers.Getlist_menu)
			clientmenu.GET("/getparent", developers.Getparent_menu)
			clientmenu.POST("/add", developers.Add_menu)
			clientmenu.POST("/lock", developers.Add_lock)
			clientmenu.DELETE("/del", developers.Del_menu)
			clientmenu.POST("/isbase", developers.IsBase)
		}
		//测试附件
		testattachment := developer.Group("/testattachment")
		{
			testattachment.POST("/add", developers.Add_cate_testattachment)
			testattachment.GET("/getparentcate", developers.Getparent_cate_testattachment)
			testattachment.GET("/getcate", developers.Get_cate_testattachment)
			testattachment.DELETE("/delcate", developers.Del_cate_testattachment)
			testattachment.POST("/upcate", developers.Up_cate_testattachment)
			testattachment.GET("/getlist", developers.Getlist_testattachment)
			testattachment.POST("/addcontent", developers.Add_content_testattachment)
			testattachment.DELETE("/delcontent", developers.Del_content_testattachment)
			testattachment.POST("/upcontent", developers.Up_content_testattachment)
		}
		//测试附件
		visiturl := developer.Group("/visiturl")
		{
			visiturl.POST("/add", developers.Add_cate_visiturl)
			visiturl.GET("/getparentcate", developers.Getparent_cate_visiturl)
			visiturl.GET("/getcate", developers.Get_cate_visiturl)
			visiturl.DELETE("/delcate", developers.Del_cate_visiturl)
			visiturl.POST("/upcate", developers.Up_cate_visiturl)
			visiturl.GET("/getlist", developers.Getlist_visiturl)
			visiturl.POST("/addcontent", developers.Add_content_visiturl)
			visiturl.DELETE("/delcontent", developers.Del_content_visiturl)
			visiturl.POST("/upcontent", developers.Up_content_visiturl)
		}
	}
	//操作cmd命令
	cmdshell_r := R.Group("/cmdshell")
	{
		//windows系统
		cmds := cmdshell_r.Group("/cmd")
		{
			cmds.POST("/reloadnginx", cmdshell.Reloadnginx)
		}
	}
}
