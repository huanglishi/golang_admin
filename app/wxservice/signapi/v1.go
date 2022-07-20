package signapi

import (
	"basegin/utils/Toolconf"
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"basegin/utils/wechat"
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type WxSign struct {
	// Appid 公众号appid
	Appid string
	// AppSecret 公众号秘钥
	AppSecret string
	// TokenRdsKey access_token缓存key
	TokenRdsKey string
	// TicketRdsKey ticket缓存key
	TicketRdsKey string
}

// WxJsSign
type WxJsSign struct {
	Appid     string `json:"appid"`
	Noncestr  string `json:"noncestr"`
	Timestamp string `json:"timestamp"`
	Url       string `json:"url"`
	Signature string `json:"signature"`
	Name      string `json:"name"`
	Des       string `json:"des"`
	Icon      string `json:"icon"`
}

//S-SDK使用权限签名算法signature
func JSDKsignature(context *gin.Context) {
	url := context.DefaultQuery("url", "0")
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	//获取access_token
	wechataccount_data, _ := DB().Table("merchant_wechataccount").Where("accountID", user.Accountid).Fields("id,appid,appsecret,access_token,access_token_time,subscribe_msg,name,des,icon").First()
	oldtime := wechataccount_data["access_token_time"].(int64)
	var access_token = ""
	if time.Now().Unix()-oldtime >= 7200 {
		accessToken, err := wechat.RequestToken(wechataccount_data["appid"].(string), wechataccount_data["appsecret"].(string))
		if err == nil {
			DB().Table("merchant_wechataccount").Where("id", wechataccount_data["id"]).Data(map[string]interface{}{"access_token": accessToken, "access_token_time": time.Now().Unix()}).Update()
			access_token = accessToken
		} else {
			results.Failed(context, "更新access_token失败！", err.Error())
			context.Abort()
			return
		}
	} else {
		access_token = wechataccount_data["access_token"].(string)
	}
	if access_token != "" {
		//获取ticket
		getticket, err := get_ticket(access_token, user.Accountid)
		if err != nil {
			results.Failed(context, "获取ticket失败！", err.Error())
		} else {
			if getticket["errmsg"].(string) != "ok" {
				results.Failed(context, getticket["errmsg"].(string), nil)
			} else { //ticket成功
				signature, _err := GetJsSign(url, getticket["ticket"].(string), wechataccount_data["appid"].(string))
				signature.Name = wechataccount_data["name"].(string)
				signature.Des = wechataccount_data["des"].(string)
				signature.Icon = wechataccount_data["icon"].(string)
				results.Success(context, "获取jsdk签名成功", signature, _err)
			}
		}
	} else {
		results.Failed(context, "获取access_token失败！", nil)
	}
}

//获取ticket
func get_ticket(access_token string, Accountid int64) (map[string]interface{}, error) {
	wechataccount_data, _ := DB().Table("merchant_wechataccount").Where("accountID", Accountid).Fields("id,jsapi_ticket,jsapi_ticket_time").First()
	oldtime := wechataccount_data["jsapi_ticket_time"].(int64)
	if time.Now().Unix()-oldtime >= 7200 {
		getticket, err := utils.HttpGet(Toolconf.AppConfig.String("wechat.getticket"), map[string]interface{}{"access_token": access_token, "type": "jsapi"})
		if err != nil {
			return nil, err
		} else {
			return getticket, nil
		}
	} else {
		return map[string]interface{}{"ticket": wechataccount_data["jsapi_ticket"], "errcode": 0, "errmsg": "ok", "expires_in": 7200}, nil
	}
}

// 获取signature
func GetJsSign(url string, jsTicket string, Appid string) (*WxJsSign, error) {
	// splite url
	urlSlice := strings.Split(url, "#")
	jsSign := &WxJsSign{
		Appid:     Appid,
		Noncestr:  utils.RandString(16),
		Timestamp: strconv.FormatInt(time.Now().UTC().Unix(), 10),
		Url:       urlSlice[0],
	}
	jsSign.Signature = Signature(jsTicket, jsSign.Noncestr, jsSign.Timestamp, jsSign.Url)
	return jsSign, nil
}

// Signature
func Signature(jsTicket, noncestr, timestamp, url string) string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", jsTicket, noncestr, timestamp, url)))
	return fmt.Sprintf("%x", h.Sum(nil))
}
