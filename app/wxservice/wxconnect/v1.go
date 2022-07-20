package wxconnect

import (
	"basegin/app/wxservice/wxutil"
	utils "basegin/utils/tool"
	"runtime"

	"basegin/utils/wechat"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

//接入验证
func WXCheckSignature(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")
	acid := c.Query("acid")
	Token, _ := DB().Table("merchant_wechataccount").Where("accountID", acid).Value("token")
	ok := wxutil.CheckSignature(signature, timestamp, nonce, Token.(string))
	if !ok {
		log.Println("微信公众号接入校验失败!")
		return
	}
	log.Println("微信公众号接入校验成功!")
	_, _ = c.Writer.WriteString(echostr)

}

// WXTextMsg 微信文本消息结构体
type WXTextMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	PicUrl       string
	MediaId      string
	Format       string
	Recognition  string
	ThumbMediaId string
	Location_X   string
	Location_Y   string
	Scale        int64
	Label        string
	MsgId        int64
	Event        string
	EventKey     int64
	Ticket       string
	Title        string
	Description  string
	Url          string
	Latitude     string
	Longitude    string
	Precision    string
}

//接收推送消息
func WXMsgReceive(c *gin.Context) {
	acid := c.Query("acid")
	var textMsg WXTextMsg
	err := c.ShouldBindXML(&textMsg)
	if err != nil {
		log.Printf("[消息接收] - XML数据包解析失败: %v\n", err)
		return
	}
	if textMsg.MsgType == "text" {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("<file：%s,line：%v>：\n[消息接收] 文本消息 ， 消息内容为: %s\n", file, line, textMsg.Content)
	} else if textMsg.MsgType == "image" {
		log.Printf("[消息接收]图片消息  ，图片链接（由系统生成）: %s, 图片消息媒体id: %s\n", textMsg.PicUrl, textMsg.MediaId)
	} else if textMsg.MsgType == "voice" {
		log.Printf("[消息接收] 语音消息 ，语音消息媒体id: %s, 语音格式: %s\n", textMsg.MediaId, textMsg.Format)
	} else if textMsg.MsgType == "video" {
		log.Printf("[消息接收]视频消息， 消息MediaId: %s, 消息ThumbMediaId: %s\n", textMsg.MediaId, textMsg.ThumbMediaId)
	} else if textMsg.MsgType == "location" {
		log.Printf("[消息接收] 地理位置消息， 纬度: %s, 经度: %s\n", textMsg.Location_X, textMsg.Location_Y)
	} else if textMsg.MsgType == "link" {
		log.Printf("[消息接收] 地理位置消息， 消息标题: %s, 消息描述: %s，, 消息链接: %s\n", textMsg.Title, textMsg.Description, textMsg.Url)
	} else if textMsg.MsgType == "event" {
		//接收事件推送
		if textMsg.Event == "subscribe" {
			log.Printf("[接收事件推送]  关注subscribe")
			msgtext := DoSubscribe(textMsg.FromUserName, acid)
			if msgtext != "" { //回复语不为空
				wechat.WXMsgReply(c, textMsg.ToUserName, textMsg.FromUserName, msgtext)
			}
		} else if textMsg.Event == "unsubscribe" {
			log.Printf("[接收事件推送]  取消订阅")
			DoUnsubscribe(textMsg.FromUserName)
		} else if textMsg.Event == "SCAN" {
			log.Printf("[接收事件推送]  二维码,EventKey %v, 消息Ticket: %s\n", textMsg.EventKey, textMsg.Ticket)
		} else if textMsg.Event == "LOCATION" {
			log.Printf("[接收事件推送] 上报地理位置事件,纬度: %v, 经度: %s, 位置精度: %s\n", textMsg.Latitude, textMsg.Longitude, textMsg.Precision)
		} else if textMsg.Event == "CLICK" {
			log.Printf("[接收事件推送]自定义菜单事件,事件KEY值: %v\n", textMsg.EventKey)
		} else if textMsg.Event == "VIEW" {
			log.Printf("[接收事件推送]点击菜单跳转链接时的事件推送,事件KEY值: %vs\n", textMsg.EventKey)
		}
	}
}

//用户关注公众号
func DoSubscribe(Openid string, acid string) string {
	datafield, _ := DB().Table("merchant_wechataccount").Where("accountID", acid).Fields("id,appid,appsecret,access_token,access_token_time,subscribe_msg").First()
	oldtime := datafield["access_token_time"].(int64)
	var access_token = ""
	if time.Now().Unix()-oldtime >= 7200 {
		accessToken, err := wechat.RequestToken(datafield["appid"].(string), datafield["appsecret"].(string))
		if err == nil {
			DB().Table("merchant_wechataccount").Where("id", datafield["id"]).Data(map[string]interface{}{"access_token": accessToken, "access_token_time": time.Now().Unix()}).Update()
			access_token = accessToken
		}
	} else {
		access_token = datafield["access_token"].(string)
	}
	if access_token != "" { //token有效是获取关注人信息
		urlStr := "https://api.weixin.qq.com/cgi-bin/user/info?access_token"
		parameter, err := utils.HttpGet(urlStr, map[string]interface{}{"access_token": access_token, "openid": Openid, "lang": "zh_CN"})
		if err == nil {
			//保存用户数据到本地数据库
			_wuid, _ := DB().Table("merchant_wechat_user").Where("openid", Openid).Value("wuid")
			savedata := map[string]interface{}{"wechatid": acid, "accountID": acid, "subscribe": parameter["subscribe"], "openid": Openid, "nickname": parameter["nickname"], "remark": parameter["remark"], "sex": parameter["sex"], "language": parameter["language"], "address": parameter["address"], "area": parameter["area"], "city": parameter["city"], "province": parameter["province"], "country": parameter["country"], "headimgurl": parameter["headimgurl"], "subscribe_time": parameter["subscribe_time"], "unionid": parameter["unionid"], "subscribe_scene": parameter["subscribe_scene"]}
			if _wuid != nil {
				DB().Table("merchant_wechat_user").Where("wuid", _wuid).Data(savedata).Update()
			} else {
				DB().Table("merchant_wechat_user").Data(savedata).Insert()
			}
		}
	}
	var rmsg = ""
	if datafield["subscribe_msg"] != nil {
		rmsg = datafield["subscribe_msg"].(string)
	}
	return rmsg
}

//用户取消订阅公众号
func DoUnsubscribe(openid string) {
	_wuid, _ := DB().Table("merchant_wechat_user").Where("openid", openid).Value("wuid")
	if _wuid != nil {
		DB().Table("merchant_wechat_user").Where("wuid", _wuid).Data(map[string]interface{}{"subscribe": 0, "unsubscribe_time": time.Now().Unix()}).Update()
	}
}
