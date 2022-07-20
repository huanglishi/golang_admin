package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)

//发送GET请求
func HttpGet(url_text string, data map[string]interface{}) (map[string]interface{}, error) {
	u, err := url.Parse(url_text)
	if err != nil {
		log.Fatal(err)
	}
	paras := &url.Values{}
	//设置请求参数
	for k, v := range data {
		paras.Set(k, v.(string))
	}
	u.RawQuery = paras.Encode()
	resp, err := http.Get(u.String())
	//关闭资源
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.New("request token err :" + err.Error())
	}
	jMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jMap)
	if err != nil {
		return nil, errors.New("request token response json parse err :" + err.Error())
	} else {
		return jMap, nil
	}

}

//发送POST请求
func HttpPost(url_text string, urldata map[string]interface{}, postdata map[string]interface{}, contentType string) (map[string]interface{}, error) {
	u, err := url.Parse(url_text)
	if err != nil {
		log.Fatal(err)
	}
	paras := &url.Values{}
	//设置请求参数
	for k, v := range urldata {
		paras.Set(k, v.(string))
	}
	u.RawQuery = paras.Encode()
	//json序列化
	jsonData := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(jsonData)
	jsonEncoder.SetEscapeHTML(false)
	if err := jsonEncoder.Encode(postdata); err != nil {
		return nil, errors.New("请求错误 :" + err.Error())
	}
	body := bytes.NewBufferString(string(jsonData.Bytes()))
	resp, erro := http.Post(u.String(), contentType, body)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if erro != nil {
		return nil, errors.New("请求错误 :" + erro.Error())
	}
	jMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jMap)
	if err != nil {
		return nil, errors.New(" 返回结果解析错误 :" + err.Error())
	} else {
		return jMap, nil
	}

}

//发送POST请求-备用
func HttpPost_c(url_text string, urldata map[string]interface{}, postdata map[string]interface{}, contentType string) (map[string]interface{}, error) {
	u, err := url.Parse(url_text)
	if err != nil {
		log.Fatal(err)
	}
	paras := &url.Values{}
	//设置请求参数
	for k, v := range urldata {
		paras.Set(k, v.(string))
	}
	u.RawQuery = paras.Encode()
	jsonStr, _ := json.Marshal(postdata)
	body := bytes.NewBuffer([]byte(jsonStr))
	resp, erro := http.Post(u.String(), contentType, body)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if erro != nil {
		return nil, errors.New("请求错误 :" + erro.Error())
	}
	jMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jMap)
	if err != nil {
		return nil, errors.New(" 返回结果解析错误 :" + err.Error())
	} else {
		return jMap, nil
	}

}
