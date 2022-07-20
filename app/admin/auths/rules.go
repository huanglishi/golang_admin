package auths

import (
	"basegin/app/model"
	"basegin/utils/results"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

func DB() gorose.IOrm {
	return model.DB.NewOrm()
}

//获取菜单
func Getlist(context *gin.Context) {
	rule, _ := DB().Table("admin_auth_rule").Order("weigh asc").Get()
	rulenum := GetTreeArray(rule, 0, "")
	// rulenum := make([]interface{}, 0)
	list_text := getTreeList_txt(rulenum, "title")
	results.Success(context, "查找获取菜单成功！", list_text, nil)
}

//获取菜单父级数据
func Getparent(context *gin.Context) {
	rule, _ := DB().Table("admin_auth_rule").Where("type", "menu").Fields("id,pid,title,weigh").Order("weigh asc").Get()
	rulenum := GetTreeArray(rule, 0, "")
	// rulenum := make([]interface{}, 0)
	list_text := getTreeList_txt(rulenum, "title")
	results.Success(context, "菜单父级数据！", list_text, nil)
}

//添加菜单
func Add(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	f_id := parameter["id"].(float64)
	if f_id == 0 {
		addId, err := DB().Table("admin_auth_rule").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加菜单失败", err)
		} else {
			if addId != 0 {
				DB().Table("admin_auth_rule").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("admin_auth_rule").
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
	res2, err := DB().Table("admin_auth_rule").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
func Del(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("admin_auth_rule").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除菜单失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

//tool-获取树状数组
func GetTreeArray(num []gorose.Data, pid int64, itemprefix string) []gorose.Data {
	childs := ToolFar(num, pid) //获取pid下的所有数据
	var chridnum []gorose.Data
	if childs != nil {
		var number int = 1
		var total int = len(childs)
		for _, v := range childs {
			j := ""
			k := ""
			if number == total {
				j += "└"
				k = ""
				if itemprefix != "" {
					k = "&nbsp;"
				}

			} else {
				j += "├"
				k = ""
				if itemprefix != "" {
					k = "│"
				}
			}
			spacer := ""
			if itemprefix != "" {
				spacer = itemprefix + j
			}
			v["spacer"] = spacer
			v["childlist"] = GetTreeArray(num, v["id"].(int64), itemprefix+k+"&nbsp;")
			chridnum = append(chridnum, v)
			number++
		}
	}
	return chridnum
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

//base_tool-获取pid下所有数组
func ToolFar(data []gorose.Data, pid int64) []gorose.Data {
	var mapString []gorose.Data
	for _, v := range data {
		if v["pid"].(int64) == pid {
			mapString = append(mapString, v)
		}
	}
	return mapString
}

//数组拼接
func ArrayMerge(ss ...[]gorose.Data) []gorose.Data {
	n := 0
	for _, v := range ss {
		n += len(v)
	}
	s := make([]gorose.Data, 0, n)
	for _, v := range ss {
		s = append(s, v...)
	}
	return s
}

//三元表达式、三目运算 f(2>3, "大于", false)
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
