package wxhome

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"time"

	"github.com/gin-gonic/gin"
)

//获取首页上部分页面
func Getbasetop_v1(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	//轮播数据
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	swipelist, _ := DB().Table("m_swipe_ads").Where("accountID", user.Accountid).Where("status", 0).Where("valid_time", ">=", startTime).Fields("id,title,image,link,linktype,isskip").Order("weigh asc,id desc").Get()
	//套餐
	packlist, _ := DB().Table("m_shopping_goods_content").Where("accountID", user.Accountid).Where("status", 0).Where("type", "package").Where("flag", "like", "%recom%").Fields("id,title,price,bgimg").Order("weigh asc,id desc").Get()
	//医生生
	doctorlist, _ := DB().Table("m_workermanage_staff_user").Where("accountID", user.Accountid).Where("status", 0).Where("flag", "like", "%recom%").Fields("id,name,position,professional,headimgurl").Order("weigh asc,id desc").Get()
	//产品
	goodslist, _ := DB().Table("m_shopping_goods_content").Where("accountID", user.Accountid).Where("status", 0).Where("type", "goods").Where("flag", "like", "%recom%").Fields("id,title,price,original_price,sales,thumb").Order("weigh asc,id desc").Get()
	//科普文章
	articlelist, _ := DB().Table("m_marticle_content").Where("accountID", user.Accountid).Where("status", 0).Where("flag", "like", "%recom%").Fields("id,title,thumb,des").Order("weigh asc,id desc").Get()
	results.Success(context, "获取首页上部分内容", map[string]interface{}{"swipe": swipelist, "pack": packlist, "doctor": doctorlist, "goods": goodslist, "article": articlelist}, startTime)
}
