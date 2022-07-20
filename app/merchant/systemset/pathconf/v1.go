package pathconf

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//获取内容数据
func Getlist(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("admin_dev_visiturl_domain").Where("accountID", user.Accountid).Order("id asc").Get()
	for _, val := range list {
		var visiturl_ids []interface{}
		json.Unmarshal([]byte(val["visiturl_ids"].(string)), &visiturl_ids)
		val["visiturl_ids"] = visiturl_ids
	}
	//获取选择地址数据
	//购买的权限
	choosepack, _ := DB().Table("admin_business_choosepack").Where("validtime", ">", time.Now().Unix()).Where("ispay", 1).Where("accountID", user.Accountid).Pluck("vid")
	choosepack_or, _ := DB().Table("admin_business_choosepack").Where("validtime", 0).Where("ispay", 1).Where("accountID", user.Accountid).Pluck("vid")
	choosepack_arr := choosepack.([]interface{}) //转换
	for _, val := range choosepack_or.([]interface{}) {
		choosepack_arr = append(choosepack_arr, val)
	}
	visiturl, _ := DB().Table("admin_dev_module_version").WhereIn("id", choosepack_arr).Pluck("visiturl_ids")
	var rule_ids_arr []string
	for _, v := range visiturl.([]interface{}) {
		if v != nil && v != "" {
			ids_arr := strings.Split(v.(string), `,`)
			rule_ids_arr = append(rule_ids_arr, ids_arr...)
		}
	}
	//将[] string转换为[] interface {}
	urlids := make([]interface{}, len(rule_ids_arr))
	for i, v := range rule_ids_arr {
		urlids[i] = v
	}
	urldata, _ := DB().Table("admin_dev_visiturl_content").WhereIn("id", urlids).Order("id asc").Get()
	for _, val := range urldata {
		catename, _ := DB().Table("admin_dev_visiturl_cate").Where("id", val["cid"]).Value("name")
		val["catename"] = catename
	}
	results.Success(context, "获取列表", map[string]interface{}{"list": list, "urldata": urldata}, nil)
}

//添加-更新
func Updata_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	f_id := parameter["id"].(float64)
	//数组转字符串
	visiturl_ids, _ := utils.JsonMarshalNoSetEscapeHTML(parameter["visiturl_ids"])
	parameter["visiturl_ids"] = visiturl_ids
	if f_id == 0 {
		parameter["accountID"] = user.Accountid
		addId, err := DB().Table("admin_dev_visiturl_domain").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("admin_dev_visiturl_domain").
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
func Delconf(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["id"]
	res2, err := DB().Table("admin_dev_visiturl_domain").Where("id", ids).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}
