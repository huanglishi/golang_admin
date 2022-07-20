package developers

import (
	"basegin/utils/results"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

//获取菜单
func Getlist_menu(context *gin.Context) {
	rule, _ := DB().Table("merchant_client_menu").Order("weigh asc").Get()
	rulenum := GetTreeArray(rule, 0, "")
	// rulenum := make([]interface{}, 0)
	list_text := getTreeList_txt(rulenum, "title")
	results.Success(context, "查找获取菜单成功！", list_text, nil)
}

//获取菜单父级数据
func Getparent_menu(context *gin.Context) {
	rule, _ := DB().Table("merchant_client_menu").Where("type", "menu").Fields("id,pid,title,weigh").Order("weigh asc").Get()
	rulenum := GetTreeArray(rule, 0, "")
	// rulenum := make([]interface{}, 0)
	list_text := getTreeList_txt(rulenum, "title")
	results.Success(context, "菜单父级数据！", list_text, nil)
}

//添加菜单
func Add_menu(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	f_id := parameter["id"].(float64)
	if f_id == 0 {
		addId, err := DB().Table("merchant_client_menu").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加菜单失败", err)
		} else {
			if addId != 0 {
				DB().Table("merchant_client_menu").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("merchant_client_menu").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新菜单失败", err)
		} else {
			results.Success(context, "更新成功！", res, nil)
		}
	}
}

// 更新状态
func Add_lock(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("merchant_client_menu").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		msg := "更新成功！"
		if res2 == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, res2, nil)
	}

}

// 更新基础菜单状态
func IsBase(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("merchant_client_menu").WhereIn("id", ids_arr).Data(map[string]interface{}{"isbase": parameter["status"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		msg := "更新成功！"
		if res2 == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, res2, nil)
	}

}

//删除菜单
func Del_menu(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("merchant_client_menu").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除菜单失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

//2.将getTreeArray的结果返回为二维数组
func getTreeList_txt(data []gorose.Data, field string) []gorose.Data {
	var midleArr []gorose.Data
	for _, v := range data {
		var childlist []gorose.Data
		if _, ok := v["childlist"]; ok {
			childlist = v["childlist"].([]gorose.Data)
		} else {
			childlist = make([]gorose.Data, 0)
		}
		delete(v, "childlist")
		v[field+"_txt"] = v["spacer"].(string) + " " + v[field+""].(string)
		if len(childlist) > 0 {
			v["haschild"] = 1
		} else {
			v["haschild"] = 0
		}
		if _, ok := v["id"]; ok {
			midleArr = append(midleArr, v)
		}
		if len(childlist) > 0 {
			newarr := getTreeList_txt(childlist, field)
			midleArr = ArrayMerge(midleArr, newarr)
		}
	}
	return midleArr
}

//三元表达式、三目运算 f(2>3, "大于", false)
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
