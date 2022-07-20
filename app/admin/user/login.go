package user

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//登录
func Lonin(context *gin.Context) {
	// username := context.PostForm("username")
	// password := context.PostForm("password")
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	username := parameter["username"].(string)
	password := parameter["password"].(string)
	if username == "" {
		results.Failed(context, "请提交用户账号！", nil)
		return
	}
	res, err := DB().Table("admin_user").Fields("id,password,salt,name").Where("username", username).First()
	if res == nil || err != nil {
		results.Failed(context, "账号不存在！", nil)
		return
	}
	pass := utils.Md5(password + res["salt"].(string))
	if pass != res["password"] {
		results.Failed(context, "您输入的密码不正确！", password)
		return
	}
	//token
	token := utils.GenerateToken(&utils.UserClaims{
		ID:             res["id"].(int64),
		Accountid:      0,
		Name:           res["name"].(string),
		Username:       username,
		StandardClaims: jwt.StandardClaims{},
	})
	// log.Printf("TOken解析: %v\n", token)
	DB().Table("admin_user").Where("id", res["id"]).Data(map[string]interface{}{"status": 1, "lastLoginTime": time.Now().Unix(), "lastLoginIp": utils.GetRequestIP(context)}).Update()
	results.Success(context, "登录成功！", res, token)
	// strtoken := string(token)
	// tokenstu := utils.ParseToken(strtoken)
	// log.Printf("TOken解析: %v\n", tokenstu.Name)
}

//退出登录
func Logout(context *gin.Context) {
	getuser, ok := context.Get("user") //取值 实现了跨中间件取值
	if !ok {
		results.Failed(context, "用户id不存在！", ok)
		return
	}
	user := getuser.(*utils.UserClaims)
	res, err := DB().Table("admin_user").Where("id", user.ID).Data(map[string]interface{}{"status": 0}).Update()
	if err != nil {
		results.Failed(context, "退出登录失败！", err)
	} else {
		results.Success(context, "退出登录成功！", res, nil)
	}
}

//刷新token
func Refreshtoken(context *gin.Context) {
	// token := context.PostForm("token")
	//获取post传过来的data
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
