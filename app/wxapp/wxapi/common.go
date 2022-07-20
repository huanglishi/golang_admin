package wxapi

import (
	"basegin/app/model"
	"basegin/utils/Toolconf"
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

func DB() gorose.IOrm {
	return model.DB.NewOrm()
}

/**公共方法*/
//微信登录
//获取微信账号
func Getwechataccount_v1(context *gin.Context) {
	accountID := context.DefaultQuery("aid", "0")
	data, _ := DB().Table("merchant_wechataccount").Where("accountID", accountID).Fields("id,appid,appsecret,name,qrcode").First()
	results.Success(context, "获取微信账号", data, nil)
}

//换取网页授权access_token-获取登录人信息
func Getaccesstoken_v1(context *gin.Context) {
	appid := context.DefaultQuery("appid", "")
	secret := context.DefaultQuery("secret", "")
	code := context.DefaultQuery("code", "")
	parameter, err := utils.HttpGet(Toolconf.AppConfig.String("wechat.access_token"), map[string]interface{}{"appid": appid, "secret": secret, "code": code, "grant_type": "authorization_code"})
	if err == nil {
		userinfo, _ := DB().Table("merchant_wechat_user").Where("openid", parameter["openid"]).Fields("wuid,accountID,mobile,openid,nickname,remark,sex,headimgurl").First()
		//token
		if userinfo != nil {
			token := utils.GenerateToken(&utils.UserClaims{
				ID:             userinfo["wuid"].(int64),
				Accountid:      userinfo["accountID"].(int64),
				Name:           userinfo["remark"].(string),
				Username:       userinfo["headimgurl"].(string),
				StandardClaims: jwt.StandardClaims{},
			})
			results.Success(context, "获取微信用户信息成功", userinfo, token)
		} else { //数据库无用户数据
			wx_userinfo, uerr := utils.HttpGet(Toolconf.AppConfig.String("wechat.userinfo"), map[string]interface{}{"access_token": parameter["access_token"], "openid": parameter["openid"], "lang": "zh_CN"})
			if uerr == nil {
				wxaccount, _ := DB().Table("merchant_wechataccount").Where("appid", appid).Fields("accountID").First()
				if wx_userinfo["nickname"] == nil {
					wx_userinfo["nickname"] = "游客"
				}
				if wx_userinfo["headimgurl"] == nil {
					wx_userinfo["headimgurl"] = ""
				}
				if wx_userinfo["sex"] == nil {
					wx_userinfo["sex"] = 3
					wx_userinfo["province"] = " "
					wx_userinfo["city"] = " "
					wx_userinfo["country"] = " "
				}
				var userinfo_data = map[string]interface{}{"wechatid": wxaccount["accountID"], "accountID": wxaccount["accountID"], "openid": parameter["openid"], "nickname": wx_userinfo["nickname"], "subscribe": 3, "headimgurl": wx_userinfo["headimgurl"], "sex": wx_userinfo["sex"], "province": wx_userinfo["province"], "city": wx_userinfo["city"], "country": wx_userinfo["country"]}
				wuid, _ := DB().Table("merchant_wechat_user").Data(userinfo_data).InsertGetId()
				token := utils.GenerateToken(&utils.UserClaims{
					ID:             wuid,
					Accountid:      userinfo_data["accountID"].(int64),
					Name:           wx_userinfo["nickname"].(string),
					Username:       userinfo_data["headimgurl"].(string),
					StandardClaims: jwt.StandardClaims{},
				})
				userinfo_data["wuid"] = wuid
				results.Success(context, "获取微信用户信息成功", userinfo_data, token)
			} else {
				results.Failed(context, "拉取用户信息失败", uerr)
			}
		}
	} else {
		results.Failed(context, "换取网页授权失败", err)
	}
}

//刷新token
func Refreshtoken_v1(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	token := parameter["token"].(string)
	newtoken := utils.Refresh(token)
	context.JSON(200, gin.H{
		"code":  0,
		"token": newtoken,
		"msg":   "刷新token",
		"data":  nil,
		"time":  time.Now().Unix(),
	})
}

//常规登录-密码账号
func Commonlogin_v1(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	userinfo, _ := DB().Table("merchant_wechat_user").Where("username", parameter["username"]).OrWhere("mobile", parameter["username"]).Fields("wuid,accountID,mobile,openid,nickname,remark,sex,headimgurl,password").First()
	if userinfo == nil {
		results.Failed(context, "账号不存在！", nil)
	} else {
		password := userinfo["password"].(string)
		pass := utils.Md5(parameter["password"].(string))
		if pass != password {
			results.Failed(context, "您输入的密码不正确！", nil)
		} else { //验证通过
			//token
			token := utils.GenerateToken(&utils.UserClaims{
				ID:             userinfo["wuid"].(int64),
				Accountid:      userinfo["accountID"].(int64),
				Name:           userinfo["remark"].(string),
				Username:       userinfo["headimgurl"].(string),
				StandardClaims: jwt.StandardClaims{},
			})
			results.Success(context, "登录成功获取用户信息", userinfo, token)
		}
	}
}
