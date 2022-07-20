package wxcommon

import (
	"basegin/utils/Toolconf"
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"basegin/utils/wechat"
	"time"

	"github.com/gin-gonic/gin"
)

//获取绑定的微信二维码
func GetQrcode_v1(context *gin.Context) {
	//获取post传过来的data
	codekey := context.DefaultQuery("codekey", "")
	action_name := context.DefaultQuery("action_name", "QR_STR_SCENE")
	expire_seconds := context.DefaultQuery("expire_seconds", "60")
	//取值用户信息
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	wechataccount, _ := DB().Table("merchant_wechataccount").Where("accountID", user.Accountid).Fields("id,name,appid,appsecret,access_token,access_token_time,subscribe_msg").First()
	oldtime := wechataccount["access_token_time"].(int64)
	var access_token = ""
	if time.Now().Unix()-oldtime >= 7200 {
		accessToken, err := wechat.RequestToken(wechataccount["appid"].(string), wechataccount["appsecret"].(string))
		if err != nil {
			results.Failed(context, "获取accessToken失败", err.Error())
			context.Abort()
			return
		} else {
			DB().Table("merchant_wechataccount").Where("id", wechataccount["id"]).Data(map[string]interface{}{"access_token": accessToken, "access_token_time": time.Now().Unix()}).Update()
			access_token = accessToken
		}
	} else {
		access_token = wechataccount["access_token"].(string)
	}
	if access_token != "" { //token有效是
		var postdata map[string]interface{}
		if action_name == "QR_SCENE" {
			postdata = map[string]interface{}{"expire_seconds": expire_seconds, "action_name": action_name, "action_info": map[string]interface{}{"scene": map[string]interface{}{"scene_id": codekey}}}
		} else if action_name == "QR_STR_SCENE" {
			postdata = map[string]interface{}{"expire_seconds": expire_seconds, "action_name": action_name, "action_info": map[string]interface{}{"scene": map[string]interface{}{"scene_str": codekey}}}
		} else if action_name == "QR_LIMIT_SCENE" {
			postdata = map[string]interface{}{"action_name": action_name, "action_info": map[string]interface{}{"scene": map[string]interface{}{"scene_id": codekey}}}
		} else if action_name == "QR_LIMIT_STR_SCENE" {
			postdata = map[string]interface{}{"action_name": action_name, "action_info": map[string]interface{}{"scene": map[string]interface{}{"scene_str": codekey}}}
		}
		res, err := utils.HttpPost(Toolconf.AppConfig.String("wechat.get_qrcode"), map[string]interface{}{"access_token": access_token}, postdata, "application/json; encoding=utf-8")
		if err != nil {
			results.Failed(context, err.Error(), nil)
		} else {
			results.Success(context, "获取绑定的微信二维码成功！", res, nil)
		}
	} else {
		results.Failed(context, "accessToken无效", nil)
	}

	context.Abort()
	return
}
