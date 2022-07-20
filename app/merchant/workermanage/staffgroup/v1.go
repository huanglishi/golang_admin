package staffgroup

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

//添加
func Addcate_v1(context *gin.Context) {
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
		addId, err := DB().Table("m_workermanage_staff_group").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			if addId != 0 {
				DB().Table("m_workermanage_staff_group").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("m_workermanage_staff_group").
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
func Getcate_v1(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("m_workermanage_staff_group").Where("accountID", user.Accountid).Order("weigh asc,id asc").Get()
	for k, val := range list {
		child_ids, _ := DB().Table("m_workermanage_staff_group").Where("pid", val["id"]).Pluck("id")
		cate_ids := child_ids.([]interface{})
		cate_ids = append(cate_ids, val["id"])
		getcount, _ := DB().Table("m_workermanage_staff_group_bind").WhereIn("group_id", cate_ids).Count()
		val["count"] = getcount
		if k == 0 {
			totalcount, _ := DB().Table("m_workermanage_staff_user").Where("accountID", user.Accountid).Count()
			val["total"] = totalcount
		}
	}
	list_tree := GetTreeArray_only(list, 0)
	results.Success(context, "获取分类列表", list_tree, nil)
}

//获取父级数据
func Getparentcate_v1(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("m_workermanage_staff_group").Where("status", 0).Where("accountID", user.Accountid).Fields("id,pid,name").Order("id asc").Get()
	rulenum := GetTreeArray(list, 0, "")
	list_text := GetTreeList_txt(rulenum, "name")
	results.Success(context, "获取父级数据", list_text, nil)
}

//获取分类树
func Getcatetree_v1(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("m_workermanage_staff_group").Where("status", 0).Where("accountID", user.Accountid).Fields("id,pid,name").Order("id asc").Get()
	menu_tree := GetTreeArray_only(list, 0)
	results.Success(context, "获取分类树", menu_tree, nil)
}

// 更新分类数据
func Upcate_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("m_workermanage_staff_group").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
func Delcate_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("m_workermanage_staff_group").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}
