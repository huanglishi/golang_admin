package wxsetting

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	. "basegin/utils/wechat"
	"encoding/json"
	"io/ioutil"
	"log"
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
		addId, err := DB().Table("merchant_wechataccount").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "账号设置成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("merchant_wechataccount").
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

//获取微信账号信息
func GetInfo(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	datafield, _ := DB().Table("merchant_wechataccount").Where("accountID", user.Accountid).Fields("id,name,wxaccount,appid,original,appsecret,qrcode,icon,subscribe_msg,token,encodingaeskey,des").First()
	results.Success(context, "获取微信账号信息", datafield, user.Accountid)
}

//更新字段数据
func Upfield(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res, err := DB().Table("merchant_wechataccount").Where("id", parameter["id"]).Data(map[string]interface{}{parameter["field"].(string): parameter["val"]}).Update()
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

//检查连接
func Checkapi(context *gin.Context) {
	getuser, _ := context.Get("user") //取值
	user := getuser.(*utils.UserClaims)
	datafield, _ := DB().Table("merchant_wechataccount").Where("accountID", user.Accountid).Fields("id,appid,appsecret").First()
	accessToken, err := RequestToken(datafield["appid"].(string), datafield["appsecret"].(string))
	if err != nil {
		log.Println(err)
		results.Failed(context, "检查连接失败！", err.Error())
	} else {
		DB().Table("merchant_wechataccount").Where("id", datafield["id"]).Data(map[string]interface{}{"access_token": accessToken, "access_token_time": time.Now().Unix()}).Update()
		results.Success(context, "检查连接成功", accessToken, nil)
	}
}
