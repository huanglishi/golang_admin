package wechat

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"

	"github.com/gin-gonic/gin"
)

//获取获取表情
func GetemojiDB_v1(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	//获取店铺信息
	setting, _err := DB().Table("m_shopping_setting").Where("accountID", user.Accountid).Fields("id,name,tel,address,location,location_name,logo,door_img,notice_msg,des").First()
	//获取分类
	catedata, _err := DB().Table("m_shopping_goods_cate").Where("accountID", user.Accountid).Where("status", 0).Where("show_nav", 1).Fields("id,name").Order("weigh asc,id desc").Get()
	results.Success(context, "获取基础数据", map[string]interface{}{"setting": setting, "catedata": catedata}, _err)
}
