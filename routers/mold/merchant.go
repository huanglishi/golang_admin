package rmold

import (
	"basegin/app/merchant/auths"
	"basegin/app/merchant/material/marticlecate"
	"basegin/app/merchant/material/marticlecontent"
	"basegin/app/merchant/shopping/appointment"
	"basegin/app/merchant/shopping/goodscate"
	"basegin/app/merchant/shopping/goodscontent"
	"basegin/app/merchant/shopping/goodsorder"
	"basegin/app/merchant/shopping/mswipe"
	"basegin/app/merchant/shopping/packagecontent"
	"basegin/app/merchant/shopping/packageorder"
	"basegin/app/merchant/shopping/shoppingsetting"
	"basegin/app/merchant/systemset/pathconf"
	"basegin/app/merchant/user"
	"basegin/app/merchant/website/edithtmltpl"
	"basegin/app/merchant/website/friendlylink"
	"basegin/app/merchant/website/leavemessage"
	"basegin/app/merchant/website/sitearticle"
	"basegin/app/merchant/website/sitecate"
	"basegin/app/merchant/website/sitemodule"
	"basegin/app/merchant/website/sites"
	"basegin/app/merchant/wechat/material"
	"basegin/app/merchant/wechat/wxcommon"
	"basegin/app/merchant/wechat/wxmenu"
	"basegin/app/merchant/wechat/wxsetting"
	"basegin/app/merchant/wechat/wxuser"
	"basegin/app/merchant/workermanage/schedul"
	"basegin/app/merchant/workermanage/staffgroup"
	"basegin/app/merchant/workermanage/staffuser"

	"github.com/gin-gonic/gin"
)

//客户端路由
func Apim(R *gin.Engine) {

	//登录路由
	merchant := R.Group("/merchant/user")
	{
		merchant.GET("/info", user.GetInfo)
		merchant.POST("/login", user.Lonin)
		merchant.POST("/refreshtoken", user.Refreshtoken)
		merchant.POST("/logout", user.Logout)

	}
	//分组
	group := R.Group("/merchant/group")
	{
		group.GET("/getlist", auths.Getlist_group)
		group.GET("/getparent", auths.Getparent_group)
		group.GET("/grouptree", auths.Getgroup_tree)
		group.POST("/add", auths.Add_group)
		group.DELETE("/del", auths.Del_group)
	}
	//账号管理
	account := R.Group("/merchant/account")
	{
		account.GET("/getdata", auths.Get_account)
		account.POST("/add", auths.Add_account)
		account.DELETE("/del", auths.Del_account)
		account.POST("/update", auths.Up_account)
	}
	//微信管理
	wechat := R.Group("/wechat")
	{
		wxsetting_v1 := wechat.Group("/wxsetting/v1")
		{
			wxsetting_v1.POST("/submitdata", wxsetting.Submitdata)
			wxsetting_v1.GET("/getinfo", wxsetting.GetInfo)
			wxsetting_v1.POST("/upfield", wxsetting.Upfield)
			wxsetting_v1.GET("/checkapi", wxsetting.Checkapi)
		}
		wxuser_v1 := wechat.Group("/wxuser/v1")
		{
			wxuser_v1.GET("/getlist", wxuser.Getlist)
			wxuser_v1.GET("/getcate", wxuser.Get_cate)
			wxuser_v1.GET("/getcatetree", wxuser.Getcatetree)
			wxuser_v1.GET("/getparentcate", wxuser.Getparent_cate)
			wxuser_v1.POST("/addcate", wxuser.Add_cate)
			wxuser_v1.POST("/upcate", wxuser.Up_cate)
			wxuser_v1.POST("/setgroups", wxuser.Setgroups)
			wxuser_v1.DELETE("/delcate", wxuser.Del_cate)
			wxuser_v1.DELETE("/removecate", wxuser.Removecate)
			wxuser_v1.POST("/upuserinfo", wxuser.Upuserinfo)
			wxuser_v1.POST("/upstatus", wxuser.Upstatus)
			wxuser_v1.POST("/upwxdata", wxuser.Upwxdata)
			wxuser_v1.POST("/dwnuserinfo", wxuser.Dwnuserinfo)
		}
		wxmenu_v1 := wechat.Group("/wxmenu/v1")
		{
			wxmenu_v1.GET("/getmenu", wxmenu.Getmenu)
			wxmenu_v1.GET("/getmenufromwx", wxmenu.Getmenufromwx)
			wxmenu_v1.POST("/createmenu", wxmenu.Createmenu)
			wxmenu_v1.GET("/getweburl", wxmenu.Getweburl)

		}
		//微信素材
		material_v1 := wechat.Group("/material/v1")
		{
			material_v1.POST("/synmaterial", material.Synmaterial)
			material_v1.GET("/getonematerial", material.Getonematerial)
		}
		//生成带参数的二维码-通用
		wxcommon_v1 := wechat.Group("/wxcommon/v1")
		{
			wxcommon_v1.GET("/getwxqrcode", wxcommon.GetQrcode_v1)
		}
	}
	//素材管理
	material := R.Group("/material")
	{
		//文章管理
		marticle_v1 := material.Group("/marticle/v1")
		{
			//分类
			marticle_v1.POST("/addcate", marticlecate.Addcate_v1)
			marticle_v1.GET("/getcate", marticlecate.Getcate_v1)
			marticle_v1.GET("/getcatetree", marticlecate.Getcatetree_v1)
			marticle_v1.GET("/getparentcate", marticlecate.Getparentcate_v1)
			marticle_v1.DELETE("/delcate", marticlecate.Delcate_v1)
			marticle_v1.POST("/upcate", marticlecate.Upcate_v1)
			//内容
			marticle_v1.GET("/getlist", marticlecontent.Getlist)
			marticle_v1.GET("/getcontentmore", marticlecontent.Getcontentmore)
			marticle_v1.POST("/pushdata", marticlecontent.Pushcontent_v1) //提交和更新
			marticle_v1.DELETE("/delcontent", marticlecontent.Delcontent_v1)
			marticle_v1.POST("/upstatus", marticlecontent.Upstatus_v1)
		}
	}
	//商城管理
	shopping := R.Group("/shopping")
	{
		//分类管理
		cate_v1 := shopping.Group("/cate/v1")
		{
			cate_v1.POST("/addcate", goodscate.Addcate_v1)
			cate_v1.GET("/getcate", goodscate.Getcate_v1)
			cate_v1.GET("/getcatetree", goodscate.Getcatetree_v1)
			cate_v1.GET("/getparentcate", goodscate.Getparentcate_v1)
			cate_v1.DELETE("/delcate", goodscate.Delcate_v1)
			cate_v1.POST("/upcate", goodscate.Upcate_v1)
		}
		//商品管理
		goods_v1 := shopping.Group("/goods/v1")
		{
			goods_v1.GET("/getlist", goodscontent.Getlist)
			goods_v1.GET("/getcontentmore", goodscontent.Getcontentmore)
			goods_v1.POST("/pushdata", goodscontent.Pushcontent_v1) //提交和更新
			goods_v1.DELETE("/delcontent", goodscontent.Delcontent_v1)
			goods_v1.POST("/upstatus", goodscontent.Upstatus_v1)
		}
		//套餐管理
		spackage_v1 := shopping.Group("/spackage/v1")
		{
			spackage_v1.GET("/getlist", packagecontent.Getlist)
			spackage_v1.GET("/getcontentmore", packagecontent.Getcontentmore)
			spackage_v1.POST("/pushdata", packagecontent.Pushcontent_v1) //提交和更新
			spackage_v1.DELETE("/delcontent", packagecontent.Delcontent_v1)
			spackage_v1.POST("/upstatus", packagecontent.Upstatus_v1)
			spackage_v1.DELETE("/delpackage", packagecontent.Delpackages_v1)
		}
		//轮播广告
		mswipe_v1 := shopping.Group("/mswipe/v1")
		{
			mswipe_v1.GET("/getlist", mswipe.Getlist)
			mswipe_v1.GET("/getcontentmore", mswipe.Getcontentmore)
			mswipe_v1.POST("/pushdata", mswipe.Pushcontent_v1) //提交和更新
			mswipe_v1.DELETE("/delcontent", mswipe.Delcontent_v1)
			mswipe_v1.POST("/upstatus", mswipe.Upstatus_v1)
			mswipe_v1.GET("/getmaterial", mswipe.Getmaterial) //获取跳转链接素材
		}
		//设置
		setting_v1 := shopping.Group("/setting/v1")
		{
			setting_v1.POST("/submitdata", shoppingsetting.Submitdata)
			setting_v1.GET("/getinfo", shoppingsetting.GetInfo)
			setting_v1.POST("/upfield", shoppingsetting.Upfield)
		}
		//产品订单
		goodsorder_v1 := shopping.Group("/goodsorder/v1")
		{
			goodsorder_v1.GET("/getlist", goodsorder.Cetlist)
		}
		//套餐产品购买订单
		packageorder_v1 := shopping.Group("/packageorder/v1")
		{
			packageorder_v1.GET("/getlist", packageorder.Cetlist)
		}
		//套餐产品预约
		appointment_v1 := shopping.Group("/appointment/v1")
		{
			appointment_v1.GET("/getlist", appointment.Cetlist)
		}
	}
	//工作管理
	workermanage := R.Group("/workermanage")
	{
		//分类管理
		group_v1 := workermanage.Group("/group/v1")
		{
			group_v1.POST("/addcate", staffgroup.Addcate_v1)
			group_v1.GET("/getcate", staffgroup.Getcate_v1)
			group_v1.GET("/getcatetree", staffgroup.Getcatetree_v1)
			group_v1.GET("/getparentcate", staffgroup.Getparentcate_v1)
			group_v1.DELETE("/delcate", staffgroup.Delcate_v1)
			group_v1.POST("/upcate", staffgroup.Upcate_v1)
		}
		//人员
		staffs_v1 := workermanage.Group("/staffs/v1")
		{
			staffs_v1.GET("/getlist", staffuser.Getlist)
			staffs_v1.GET("/getcontentmore", staffuser.Getcontentmore_v1)
			staffs_v1.POST("/pushdata", staffuser.Pushcontent_v1) //提交和更新
			staffs_v1.DELETE("/delcontent", staffuser.Delcontent_v1)
			staffs_v1.POST("/upstatus", staffuser.Upstatus_v1)
			staffs_v1.POST("/setgroups", staffuser.Setgroups_v1)
			staffs_v1.DELETE("/removegroup", staffuser.Removegroup_v1)
		}
		//排班
		schedul_v1 := workermanage.Group("/schedul/v1")
		{
			schedul_v1.POST("/pushtpl", schedul.Pushtpl_v1)
			schedul_v1.GET("/gettpl", schedul.Gettpl_v1)
			schedul_v1.POST("/pushtimeframe", schedul.Pushtimeframe_v1)
			schedul_v1.DELETE("/deltimeframe", schedul.Deltimeframe_v1)
			schedul_v1.DELETE("/deltpl", schedul.Deltpl_v1)
			schedul_v1.POST("/getschedultpl", schedul.Getschedultpl_v1)
			schedul_v1.POST("/setschedul", schedul.Setschedul_v1)
		}
	}

	//网站管理
	website := R.Group("/website")
	{
		//站点设置
		sites_v1 := website.Group("/sites/v1")
		{
			sites_v1.POST("/pushsite", sites.Pushsite_v1)
			sites_v1.GET("/getsitelist", sites.Getsitelist_v1)
			sites_v1.GET("/getsitefield", sites.Getsitefield_v1)
			sites_v1.POST("/syncdomainssl", sites.Syncdomainssl)
		}
		//分类
		cate_v1 := website.Group("/cate/v1")
		{
			//分类
			cate_v1.POST("/add", sitecate.Add_v1)
			cate_v1.GET("/getlist", sitecate.Getlist_v1)
			cate_v1.GET("/getstree", sitecate.Gettree_v1)
			cate_v1.GET("/getparent", sitecate.Getparent_v1)
			cate_v1.DELETE("/del", sitecate.Del_v1)
			cate_v1.POST("/up", sitecate.Up_v1)
		}
		//内容
		article_v1 := website.Group("/article/v1")
		{
			article_v1.GET("/getlist", sitearticle.Getlist)
			article_v1.GET("/getcontentmore", sitearticle.Getcontentmore)
			article_v1.POST("/pushdata", sitearticle.Pushcontent_v1) //提交和更新
			article_v1.DELETE("/delcontent", sitearticle.Delcontent_v1)
			article_v1.POST("/upstatus", sitearticle.Upstatus_v1)
		}
		//留言
		sitemessage_v1 := website.Group("/sitemessage/v1")
		{
			sitemessage_v1.GET("/getsitelist", leavemessage.GetSitelist_v1)
			sitemessage_v1.GET("/getlist", leavemessage.Getlist_v1)
			sitemessage_v1.GET("/getcontentmore", leavemessage.Getcontentmore_v1)
			sitemessage_v1.DELETE("/delcontent", leavemessage.Delcontent_v1)
			sitemessage_v1.POST("/upstatus", leavemessage.Upstatus_v1)
		}
		//html网站编辑
		edithtmltpl_v1 := website.Group("/edithtmltpl")
		{
			edithtmltpl_v1.GET("/getsitefile", edithtmltpl.Getsitefile)
			edithtmltpl_v1.POST("/newdir", edithtmltpl.Newdir)
			edithtmltpl_v1.POST("/newfile", edithtmltpl.Newfile)
			edithtmltpl_v1.DELETE("/delfiles", edithtmltpl.Delfiles)
			edithtmltpl_v1.GET("/readfile", edithtmltpl.Readfile)
			edithtmltpl_v1.POST("/writefile", edithtmltpl.Writefile)
		}
		//html网站编辑
		module_v1 := website.Group("/module")
		{
			module_v1.GET("/getlist", sitemodule.Getlist)
			module_v1.POST("/add", sitemodule.Add)
			module_v1.DELETE("/del", sitemodule.Del)
		}
		//友情链接
		friendlylink_v1 := website.Group("/friendlylink")
		{
			friendlylink_v1.GET("/getlist", friendlylink.Getlist)
			friendlylink_v1.POST("/pushdata", friendlylink.Pushcontent_v1) //提交和更新
			friendlylink_v1.DELETE("/delcontent", friendlylink.Delcontent_v1)
			friendlylink_v1.POST("/upstatus", friendlylink.Upstatus_v1)
		}
	}
	//系统设置
	systemset := R.Group("/systemset")
	{
		pathconf_v1 := systemset.Group("/pathconf/v1")
		{
			pathconf_v1.GET("/getlist", pathconf.Getlist)
			pathconf_v1.POST("/updata", pathconf.Updata_v1)  //提交和更新
			pathconf_v1.DELETE("/delconf", pathconf.Delconf) //删除
		}
	}
}
