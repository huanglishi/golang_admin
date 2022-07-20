package material

import (
	"basegin/utils/Toolconf"
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"basegin/utils/wechat"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

//获取微信公众号素材
func Synmaterial(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
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
		res, err := utils.HttpPost(Toolconf.AppConfig.String("wechat.get_materiallist"), map[string]interface{}{"access_token": access_token}, parameter, "application/json; encoding=utf-8")
		if err != nil {
			//保存用户数据到本地数据库
			results.Failed(context, err.Error(), nil)
		} else {
			results.Success(context, "获取微信公众号素材", res, nil)
		}
	} else {
		results.Failed(context, "accessToken无效", nil)
	}

}

//获取永久素材
func Getonematerial(context *gin.Context) {
	media_id := context.DefaultQuery("mediaId", "") //要获取的素材的media_id
	getuser, _ := context.Get("user")               //取值用户信息
	user := getuser.(*utils.UserClaims)
	wechataccount, _ := DB().Table("merchant_wechataccount").Where("accountID", user.Accountid).Fields("id,appid,appsecret,access_token,access_token_time,subscribe_msg").First()
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
		res, err := utils.HttpPost(Toolconf.AppConfig.String("wechat.get_material"), map[string]interface{}{"access_token": access_token}, map[string]interface{}{"media_id": media_id}, "application/json; encoding=utf-8")
		if err != nil {
			//保存用户数据到本地数据库
			results.Failed(context, err.Error(), nil)
		} else {
			results.Success(context, "获取永久素材", res, nil)
		}
	} else {
		results.Failed(context, "accessToken无效", nil)
	}
}
