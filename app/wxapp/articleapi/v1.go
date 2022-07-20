package articleapi

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"

	"github.com/gin-gonic/gin"
)

//获取素材分类列表
func Getarticlecate_v1(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	catedata, _err := DB().Table("m_marticle_cate").Where("accountID", user.Accountid).Where("status", 0).Where("show_nav", 1).Fields("id,name").Order("weigh asc,id desc").Get()
	if catedata != nil {
		results.Success(context, "获取分类列表", catedata, _err)
	} else {
		results.Failed(context, "分类列表失败！", _err)
	}
}

//获取素材列表
func Getarticlelist_v1(context *gin.Context) {
	cid := context.DefaultQuery("cid", "0")
	Fdb := DB().Table("m_marticle_content").Where("status", 0)
	if cid != "0" {
		Fdb = Fdb.Where("cid", cid)
	}
	list, _err := Fdb.Fields("id,title,thumb,image,des,visits,likes,createtime,author,source").Order("weigh asc,id desc").Get()
	if _err != nil {
		results.Failed(context, "获取素材列表失败！", _err)
	} else {
		results.Success(context, "获取素材列表", list, nil)
	}
}

//获取素材详情
func Getarticledetail_v1(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	contentdata, _err := DB().Table("m_marticle_content").Where("id", id).Fields("id,title,content,visits,likes,createtime,author,source").First()
	if contentdata != nil {
		DB().Table("m_marticle_content").Where("id", id).Data(map[string]interface{}{"visits": contentdata["visits"].(int64) + 1}).Update()
		results.Success(context, "获取详情", contentdata, _err)
	} else {
		results.Failed(context, "获取详情失败！", _err)
	}
}
