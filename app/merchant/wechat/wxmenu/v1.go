package wxmenu

import (
	"basegin/utils/Toolconf"
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"basegin/utils/wechat"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//获取菜单
func Getmenu(context *gin.Context) {
	getuser, _ := context.Get("user") //取值用户信息
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("merchant_wechat_menu").Where("accountID", user.Accountid).Fields("id,title,value,status").Distinct().First()
	results.Success(context, "获取菜单", list, nil)
}

//从微信上获取菜单
func Getmenufromwx(context *gin.Context) {
	getuser, _ := context.Get("user") //取值用户信息
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
		menulist, err := utils.HttpGet(Toolconf.AppConfig.String("wechat.menu_info"), map[string]interface{}{"access_token": access_token})
		if err != nil {
			//保存用户数据到本地数据库
			results.Failed(context, "查询微信菜单失败", err)
		} else {
			results.Success(context, "获取菜单", menulist, nil)
		}
	} else {
		results.Failed(context, "accessToken无效", nil)
	}
}

//同步到微信--保存并推送
func Createmenu(context *gin.Context) {
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
		res, err := utils.HttpPost(Toolconf.AppConfig.String("wechat.menu_create"), map[string]interface{}{"access_token": access_token}, parameter, "application/json; encoding=utf-8")
		if err != nil {
			//保存用户数据到本地数据库
			results.Failed(context, err.Error(), nil)
		} else {
			errcode := res["errcode"].(float64)
			int_errcode := int(errcode)
			var zero_val int = 0
			if int_errcode != zero_val {
				results.Failed(context, "推送微信菜单失败", res)
			} else {
				wxmenu_id, _ := DB().Table("merchant_wechat_menu").Where("accountID", user.Accountid).Value("id")
				menuval, _ := utils.JsonMarshalNoSetEscapeHTML(parameter["button"])
				savedata := map[string]interface{}{"accountID": user.Accountid, "title": wechataccount["name"], "value": string(menuval), "status": 0, "createtime": time.Now().Unix()}
				if wxmenu_id != nil {
					DB().Table("merchant_wechat_menu").Where("id", wxmenu_id).Data(savedata).Update()
				} else {
					DB().Table("merchant_wechat_menu").Data(savedata).Insert()
				}
				results.Success(context, "推送菜单到微信成功", res, nil)
			}
		}
	} else {
		results.Failed(context, "accessToken无效", nil)
	}
}

//获取前端访问路径
func Getweburl(context *gin.Context) {
	getuser, _ := context.Get("user") //取值用户信息
	user := getuser.(*utils.UserClaims)
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
	list, _ := DB().Table("admin_dev_visiturl_content").WhereIn("id", urlids).Order("id asc").Get()
	for _, val := range list {
		domain, _ := DB().Table("admin_dev_visiturl_domain").Where("accountID", user.Accountid).Where("visiturl_ids", "like", "%"+(Strval(val["id"]))+"%").Value("domain")
		val["domain"] = domain
	}
	results.Success(context, "获取前端访问路径", map[string]interface{}{"list": list, "Accountid": user.Accountid}, nil)
}
