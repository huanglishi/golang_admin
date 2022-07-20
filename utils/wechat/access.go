package wechat

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func RequestToken(appid, secret string) (string, error) {
	u, err := url.Parse("https://api.weixin.qq.com/cgi-bin/token")
	if err != nil {
		log.Fatal(err)
	}
	paras := &url.Values{}
	//设置请求参数
	paras.Set("appid", appid)
	paras.Set("secret", secret)
	paras.Set("grant_type", "client_credential")
	u.RawQuery = paras.Encode()
	resp, httl_err := http.Get(u.String())
	//关闭资源
	if httl_err != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if httl_err != nil {
		return "", errors.New("request token err :" + err.Error())
	} else { //请求成功
		jMap := make(map[string]interface{})
		json_err := json.NewDecoder(resp.Body).Decode(&jMap)
		if json_err != nil {
			return "", errors.New("数据json转换错误 :" + err.Error())
		} else {
			if jMap["errcode"] == nil || jMap["errcode"] == 0 {
				accessToken, _ := jMap["access_token"].(string)
				return accessToken, nil
			} else {
				//返回错误信息
				errcode := strconv.FormatFloat(jMap["errcode"].(float64), 'E', -1, 64)
				errmsg := jMap["errmsg"].(string)
				err = errors.New("errcode：" + errcode + "，errmsg：" + errmsg)
				return "", err
			}
		}
	}
}
