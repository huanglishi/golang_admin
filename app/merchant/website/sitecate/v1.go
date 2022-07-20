package sitecate

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

//添加
func Add_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	f_id := parameter["id"].(float64)
	if f_id == 0 {
		parameter["accountID"] = user.Accountid
		addId, err := DB().Table("m_website_article_cate").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			if addId != 0 {
				DB().Table("m_website_article_cate").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("m_website_article_cate").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			results.Success(context, "更新成功！", res, nil)
		}
	}
}

//获取分类列表
func Getlist_v1(context *gin.Context) {
	site_id := context.DefaultQuery("site_id", "0")
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("m_website_article_cate").Where("accountID", user.Accountid).Where("site_id", site_id).Order("weigh asc,id asc").Get()
	for k, val := range list {
		child_ids, _ := DB().Table("m_website_article_cate").Where("pid", val["id"]).Pluck("id")
		cate_ids := child_ids.([]interface{})
		cate_ids = append(cate_ids, val["id"])
		getcount, _ := DB().Table("m_website_article_content").WhereIn("cid", cate_ids).Count()
		val["count"] = getcount
		if k == 0 {
			totalcount, _ := DB().Table("m_website_article_content").Where("accountID", user.Accountid).Where("site_id", site_id).Count()
			val["total"] = totalcount
		}
		// //判断是否是单页
		// singlepage, _ := DB().Table("m_website_module").Where("id", val["module_id"]).Value("singlepage")
		// val["singlepage"] = singlepage
	}
	list_tree := GetTreeArray_only(list, 0)
	results.Success(context, "获取分类列表", list_tree, nil)
}

//获取父级数据
func Getparent_v1(context *gin.Context) {
	site_id := context.DefaultQuery("site_id", "0")
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("m_website_article_cate").Where("status", 0).Where("accountID", user.Accountid).Where("site_id", site_id).Fields("id,pid,catname").Order("id asc").Get()
	rulenum := GetTreeArray(list, 0, "")
	list_text := GetTreeList_txt(rulenum, "catname")
	modulelist, _ := DB().Table("m_website_module").Where("accountID", user.Accountid).Where("site_id", site_id).Fields("id,name").Order("id asc").Get()
	results.Success(context, "获取分类树", map[string]interface{}{"catelist": list_text, "modulelist": modulelist}, nil)
}

//获取分类树
func Gettree_v1(context *gin.Context) {
	site_id := context.DefaultQuery("site_id", "0")
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("m_website_article_cate").Where("status", 0).Where("accountID", user.Accountid).Where("site_id", site_id).Fields("id,pid,catname").Order("id asc").Get()
	menu_tree := GetTreeArray_only(list, 0)
	results.Success(context, "获取分类树", menu_tree, nil)
}

// 更新分类数据
func Up_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("m_website_article_cate").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		ty := "启用"
		var ty_int float64 = 1
		if parameter["status"].(float64) == ty_int {
			ty = "锁定"
		}
		msg := ty + "成功！"
		if res2 == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, res2, nil)
	}
}

//删除分类
func Del_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("m_website_article_cate").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败！", err)
	} else {
		DB().Table("m_website_article_content").Where("cid", parameter["id"]).Delete()
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}
