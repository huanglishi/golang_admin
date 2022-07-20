package auths

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

//获取
func Getlist_group(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	group_ids, _ := DB().Table("admin_auth_group_access").Where("uid", user.ID).Pluck("group_id")
	gruop_ids_int := Childrengroup(group_ids)
	rule, _ := DB().Table("admin_auth_group").Fields("id,pid,name,rules,checkedKeys,createtime,des,status,data_access").Order("id asc").Get()
	for _, v := range rule {
		if v["rules"] == "*" {
			menu, _ := DB().Table("admin_auth_rule").Where("status", 0).Fields("id,weigh").AddFields("pid as parentId,title as name").Order("weigh asc").Get()
			v["permissions"] = menu
		} else {
			ids_arr := strings.Split(v["rules"].(string), `,`)
			var ids_arr_int []interface{}
			ids_buy, _ := json.Marshal(&ids_arr)
			json.Unmarshal([]byte(ids_buy), &ids_arr_int)
			menu, _ := DB().Table("admin_auth_rule").WhereIn("id", ids_arr_int).Where("status", 0).Fields("id,weigh").AddFields("pid as parentId,title as name").Order("weigh asc").Get()
			v["permissions"] = menu
		}
	}
	rulenum := GetTreeArray(rule, 0, "")
	list_text := getTreeList_txt(rulenum, "name")
	var list_text_arr []interface{}
	for _, lv := range list_text {
		for _, cid := range gruop_ids_int {
			if lv["id"] == cid {
				list_text_arr = append(list_text_arr, lv)
			}
		}
	}
	results.Success(context, "查找获取成功！", list_text_arr, gruop_ids_int)
}

//获取父级数据
func Getparent_group(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	group_ids, _ := DB().Table("admin_auth_group_access").Where("uid", user.ID).Pluck("group_id")
	gruop_ids_int := Childrengroup(group_ids)
	list, _ := DB().Table("admin_auth_group").Where("status", 0).Fields("id,pid,name,rules,checkedKeys").Order("id asc").Get()
	rulenum := GetTreeArray(list, 0, "")
	list_text := getTreeList_txt(rulenum, "name")
	var parameter = map[string]interface{}{"group": nil, "menu": nil}
	var list_text_arr []interface{}
	for _, lv := range list_text {
		for _, cid := range gruop_ids_int {
			if lv["id"] == cid {
				list_text_arr = append(list_text_arr, lv)
			}
		}
	}
	parameter["group"] = list_text_arr
	//权限菜单
	//添加权限
	roles_ids, _ := DB().Table("admin_auth_group").WhereIn("id", group_ids.([]interface{})).Pluck("rules")
	var menu []gorose.Data
	if utils.In_array("*", roles_ids) {
		_menu, _ := DB().Table("admin_auth_rule").Where("status", 0).Fields("id,pid,weigh,title").Order("weigh asc").Get()
		menu = _menu
	} else {
		var rule_ids_arr []string
		for _, v := range roles_ids.([]interface{}) {
			if v != nil || v != "" {
				ids_arr := strings.Split(v.(string), `,`)
				rule_ids_arr = append(rule_ids_arr, ids_arr...)
			}
		}
		ids_buy, _ := json.Marshal(&rule_ids_arr)
		var ids_arr_int []interface{}
		json.Unmarshal([]byte(ids_buy), &ids_arr_int)
		_menu, _ := DB().Table("admin_auth_rule").Where("status", 0).WhereIn("id", ids_arr_int).Fields("id,pid,weigh,title").Order("weigh asc").Get()
		menu = _menu
	}
	for _, mv := range menu {
		mv["scopedSlots"] = map[string]interface{}{"title": "title"}
	}
	menu_tree := GetTreeArray_only(menu, 0)
	parameter["menu"] = menu_tree
	parameter["menuarr"] = menu
	results.Success(context, "父级数据！", parameter, roles_ids)
}

//获取父级数据tree
func Getgroup_tree(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	group_ids, _ := DB().Table("admin_auth_group_access").Where("uid", user.ID).Pluck("group_id")
	gruop_ids_int := Childrengroup(group_ids)
	list, _ := DB().Table("admin_auth_group").Where("status", 0).Fields("id,pid,name,rules,checkedKeys").Order("id asc").Get()
	rulenum := GetTreeArray(list, 0, "")
	list_text := getTreeList_txt(rulenum, "name")
	var list_text_arr []interface{}
	for _, lv := range list_text {
		for _, cid := range gruop_ids_int {
			if lv["id"] == cid {
				list_text_arr = append(list_text_arr, lv)
			}
		}
	}
	results.Success(context, "权限分组数据！", list_text_arr, nil)
}

//添加
func Add_group(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	parameter["createtime"] = time.Now().Unix()
	f_id := parameter["id"].(float64)
	if f_id == 0 {
		addId, err := DB().Table("admin_auth_group").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			if addId != 0 {
				DB().Table("admin_auth_group").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("admin_auth_group").
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

//删除
func Del_group(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("admin_auth_group").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

//tool-获取树状数组
func GetTreeArray_only(num []gorose.Data, pid int64) []gorose.Data {
	childs := ToolFar(num, pid) //获取pid下的所有数据
	var chridnum []gorose.Data
	if childs != nil {
		for _, v := range childs {
			v["children"] = GetTreeArray_only(num, v["id"].(int64))
			chridnum = append(chridnum, v)
		}
	}
	return chridnum
}

//获取全部子集
func Childrengroup(group_ids interface{}) []interface{} {
	var chriddis []interface{}
	grouplist, _ := DB().Table("admin_auth_group").Fields("id,pid").Order("id asc").Get()
	for _, v := range group_ids.([]interface{}) {
		pid := v.(int64)
		chriddis = append(chriddis, pid)
		childs := Getallchildren(grouplist, pid) //获取pid下的所有数据
		childs_arr := gettree_to_arr(childs)
		for _, grouparr := range childs_arr {
			chriddis = append(chriddis, grouparr["id"])
		}
	}
	return chriddis
}
func gettree_to_arr(data []gorose.Data) []gorose.Data {
	var midleArr []gorose.Data
	for _, v := range data {
		var childlist []gorose.Data
		if _, ok := v["children"]; ok {
			childlist = v["children"].([]gorose.Data)
		} else {
			childlist = make([]gorose.Data, 0)
		}
		delete(v, "children")
		if _, ok := v["id"]; ok {
			midleArr = append(midleArr, v)
		}
		if len(childlist) > 0 {
			newarr := gettree_to_arr(childlist)
			midleArr = ArrayMerge(midleArr, newarr)
		}
	}
	return midleArr
}

//获取全部子id
func Getallchildren(num []gorose.Data, pid int64) []gorose.Data {
	childs := ToolFar(num, pid) //获取pid下的所有数据
	var chridnum []gorose.Data
	if childs != nil {
		for _, v := range childs {
			v["children"] = GetTreeArray_only(num, v["id"].(int64))
			chridnum = append(chridnum, v)
		}
	}
	return chridnum
}
