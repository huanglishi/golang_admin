package shoppingsetting

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

//账号设置
func Submitdata(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	f_id := parameter["id"].(float64)
	if f_id == 0 {
		parameter["uid"] = user.ID
		parameter["accountID"] = user.Accountid
		parameter["createtime"] = time.Now().Unix()
		addId, err := DB().Table("m_shopping_setting").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "设置成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("m_shopping_setting").
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

//获取店铺信息
func GetInfo(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	datafield, _ := DB().Table("m_shopping_setting").Where("accountID", user.Accountid).Fields("id,name,tel,city,cityval,address,location,location_name,logo,door_img,notice_msg,des").First()
	results.Success(context, "获取店铺信息", datafield, user.Accountid)
}

//更新字段数据
func Upfield(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res, err := DB().Table("m_shopping_setting").Where("id", parameter["id"]).Data(map[string]interface{}{parameter["field"].(string): parameter["val"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		msg := "更新成功！"
		if res == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, res, nil)
	}
}
