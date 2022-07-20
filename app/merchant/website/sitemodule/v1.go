package sitemodule

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
)

//获取内容数据
func Getlist(context *gin.Context) {
	site_id := context.DefaultQuery("site_id", "0")
	keywords := context.DefaultQuery("keywords", "")
	_pageNo := context.DefaultQuery("pageNo", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(_pageNo)
	pageSize, _ := strconv.Atoi(_pageSize)
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	Fdb := DB().Table("m_website_module").Where("accountID", user.Accountid).Where("site_id", site_id)
	if keywords != "" {
		Fdb = Fdb.Where("name", "like", "%"+keywords+"%")
	}
	getfield := "id,accountID,uid,name,model_name,fcatdir,template_list,template_show,singlepage"
	list, err := Fdb.Fields(getfield).Limit(pageSize).Page(pageNo).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "查找失败", err)
	} else {
		// 统计数据
		var totalCount int64
		totalCount, _ = Fdb.Count()
		_pageSize := int64(pageSize)
		totalPage := totalCount / _pageSize
		results.Success(context, "查找成功！", map[string]interface{}{
			"pageNo":     pageNo,
			"pageSize":   pageSize,
			"totalCount": totalCount,
			"totalPage":  totalPage,
			"data":       list}, nil)
	}
}

//添加
func Add(context *gin.Context) {
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
		addId, err := DB().Table("m_website_module").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("m_website_module").
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

//删除内容
func Del(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("m_website_module").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}
