package sites

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

//添加站点
func Pushsite_v1(context *gin.Context) {
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
		parameter["createtime"] = time.Now().Unix()
		addId, err := DB().Table("m_website_site").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("m_website_site").
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

//获取网站列表
func Getsitelist_v1(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("m_website_site").Where("accountID", user.Accountid).Fields("id,name,thumb,step,sslcert,domain,des,status").Order("id asc").Get()
	results.Success(context, "获取网站列表", list, nil)
}

//获取指定字段数据
func Getsitefield_v1(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	fields := context.DefaultQuery("fields", "")
	data, _ := DB().Table("m_website_site").Where("id", id).Fields(fields).First()
	results.Success(context, "获取指定字段数据", data, nil)
}

//提交域名和ssl证书变更
func Syncdomainssl(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	var updata map[string]interface{}
	var one float64 = 1
	if parameter["type"].(float64) == one {
		updata = map[string]interface{}{"apply_domain": 1}
	} else {
		updata = map[string]interface{}{"apply_ssl": 1}
	}
	log.Println(updata)
	res, err := DB().Table("m_website_site").
		Data(updata).
		Where("id", parameter["id"]).
		Update()
	if err != nil {
		results.Failed(context, "提交失败", err)
	} else {
		results.Success(context, "提交成功！", res, nil)
	}
}
