package wxuser

import (
	"basegin/utils/Toolconf"
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"basegin/utils/wechat"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

//添加
func Add_cate(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	f_id := parameter["id"].(float64)
	if f_id == 0 {
		parameter["accountID"] = user.Accountid
		addId, err := DB().Table("merchant_wechat_user_tags").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("merchant_wechat_user_tags").
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

//获取分类列表
func Get_cate(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("merchant_wechat_user_tags").Where("accountID", user.Accountid).Order("id asc").Get()
	for k, val := range list {
		child_ids, _ := DB().Table("merchant_wechat_user_tags").Where("pid", val["id"]).Pluck("id")
		cate_ids := child_ids.([]interface{})
		cate_ids = append(cate_ids, val["id"])
		getcount, _ := DB().Table("merchant_wechat_user_addtags").WhereIn("tags_id", cate_ids).Count()
		val["count"] = getcount
		if k == 0 {
			totalcount, _ := DB().Table("merchant_wechat_user").Where("accountID", user.Accountid).Count()
			val["total"] = totalcount
		}
	}
	list_tree := GetTreeArray_only(list, 0)
	results.Success(context, "获取分类列表", list_tree, nil)
}

//获取父级数据
func Getparent_cate(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("merchant_wechat_user_tags").Where("status", 0).Where("accountID", user.Accountid).Fields("id,pid,name").Order("id asc").Get()
	rulenum := GetTreeArray(list, 0, "")
	list_text := GetTreeList_txt(rulenum, "name")
	results.Success(context, "获取父级数据", list_text, nil)
}

//获取分类树
func Getcatetree(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("merchant_wechat_user_tags").Where("status", 0).Where("accountID", user.Accountid).Fields("id,pid,name").Order("id asc").Get()
	menu_tree := GetTreeArray_only(list, 0)
	results.Success(context, "获取分类树", menu_tree, nil)
}

// 更新分类数据
func Up_cate(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("merchant_wechat_user_tags").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		ty := "启用"
		var ty_int float64 = 1
		if parameter["status"].(float64) == ty_int {
			ty = "锁定"
		}
		msg := ty + "成功！"
		if res2 == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, res2, nil)
	}
}

//删除分类
func Del_cate(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("merchant_wechat_user_tags").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

/**-------------------------------------------内容-----------------------------------------*/
//获取内容数据
func Getlist(context *gin.Context) {
	_cid := context.DefaultQuery("cid", "0")
	namekey := context.DefaultQuery("name", "")
	_pageNo := context.DefaultQuery("pageNo", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	cid, _ := strconv.Atoi(_cid)
	pageNo, _ := strconv.Atoi(_pageNo)
	pageSize, _ := strconv.Atoi(_pageSize)
	var list []gorose.Data
	var err error
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	Fdb := DB().Table("merchant_wechat_user").Where("accountID", user.Accountid)
	if namekey != "" {
		Fdb = Fdb.Where("name", "like", "%"+namekey+"%")
	}
	getfield := "wuid,wechatid,accountID,subscribe,openid,nickname,remark,sex,city,address,headimgurl,subscribe_time,unionid,status,subscribe_scene,unsubscribe_time,cityval,mobile,des"
	if cid == 0 {
		_list, _err := Fdb.Fields(getfield).Limit(pageSize).Page(pageNo).Order("wuid desc").Get()
		list = _list
		err = _err
	} else {
		_cate_ids, _ := DB().Table("merchant_wechat_user_tags").Where("pid", cid).Pluck("id")
		cate_ids := _cate_ids.([]interface{})
		cate_ids = append(cate_ids, cid)
		wuids, _ := DB().Table("merchant_wechat_user_addtags").WhereIn("tags_id", cate_ids).Pluck("user_id")
		_wuids := wuids.([]interface{})
		if len(_wuids) > 0 {
			_list, _err := Fdb.WhereIn("wuid", _wuids).Fields(getfield).Limit(pageSize).Page(pageNo).Order("wuid desc").Get()
			list = _list
			err = _err
		} else {
			list = make([]gorose.Data, 0, 1)
			err = nil
		}
	}
	if err != nil {
		results.Failed(context, "查找失败", err)
	} else {
		// 统计数据
		var totalCount int64
		totalCount, _ = DB().Table("merchant_wechat_user").Where("accountID", user.Accountid).Count()
		_pageSize := int64(pageSize)
		totalPage := totalCount / _pageSize
		for _, val := range list {
			catedata, _ := DB().Table("merchant_wechat_user_addtags a").LeftJoin("merchant_wechat_user_tags b on a.tags_id = b.id").Where("a.user_id", val["wuid"]).Fields("a.id,a.tags_id,b.pid,b.name").Get()
			if catedata != nil {
				val["tags"] = catedata
			} else {
				val["tags"] = make([]interface{}, 0, 1)
			}
		}
		results.Success(context, "查找成功！", map[string]interface{}{
			"pageNo":     pageNo,
			"pageSize":   pageSize,
			"totalCount": totalCount,
			"totalPage":  totalPage,
			"data":       list}, nil)
	}
}

//删除内容-批量
func Del_content(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("merchant_wechat_user").WhereIn("wuid", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 更新用户状态
func Upstatus(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("merchant_wechat_user").WhereIn("wuid", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		msg := "更新成功！"
		if res2 == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, res2, nil)
	}
}

// 获取公众号信息
func Upwxdata(context *gin.Context) {
	next_openid := context.DefaultQuery("next_openid", "") //要获取的素材的media_id
	getuser, _ := context.Get("user")                      //取值用户信息
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
		res, err := utils.HttpPost(Toolconf.AppConfig.String("wechat.userlist"), map[string]interface{}{"access_token": access_token}, map[string]interface{}{"next_openid": next_openid}, "application/json; encoding=utf-8")
		if err != nil {
			results.Failed(context, err.Error(), nil)
		} else {
			results.Success(context, "获取公众号信息", res, nil)
		}
	} else {
		results.Failed(context, "accessToken无效", nil)
	}
}

// 获取保存用户数据
func Dwnuserinfo(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
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
		var infou interface{}
		for _, val := range parameter["openids"].([]interface{}) {
			wx_userinfo, err := utils.HttpGet(Toolconf.AppConfig.String("wechat.user_info"), map[string]interface{}{"access_token": access_token, "openid": val, "lang": "zh_CN"})
			if err != nil {
				results.Failed(context, err.Error(), nil)
			} else {
				infou = wx_userinfo
				haseuser, _ := DB().Table("merchant_wechat_user").Where("accountID", user.Accountid).Where("openid", val).Fields("wuid").First()
				if haseuser == nil {
					var userinfo_data = map[string]interface{}{"wechatid": user.Accountid, "accountID": user.Accountid, "openid": val, "nickname": wx_userinfo["nickname"], "subscribe": wx_userinfo["subscribe"], "headimgurl": wx_userinfo["headimgurl"], "sex": wx_userinfo["sex"], "province": wx_userinfo["province"], "city": wx_userinfo["city"], "country": wx_userinfo["country"], "language": wx_userinfo["language"], "remark": wx_userinfo["remark"], "subscribe_time": wx_userinfo["subscribe_time"]}
					DB().Table("merchant_wechat_user").Data(userinfo_data).InsertGetId()
				}
			}
		}
		results.Success(context, "同步用户信息1", infou, nil)
	} else {
		results.Failed(context, "accessToken无效", nil)
	}
}

// 绑定分组
func Setgroups(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter != nil {
		//批量提交
		save_arr := []map[string]interface{}{}
		for _, val := range parameter["groupids"].([]interface{}) {
			choosepack, _ := DB().Table("merchant_wechat_user_addtags").Where("user_id", parameter["wuid"]).Where("tags_id", val).Fields("id").Distinct().First()
			if choosepack == nil {
				catedata, _ := DB().Table("merchant_wechat_user_tags").Where("pid", val).Fields("id").Distinct().First()
				if catedata == nil {
					marr := map[string]interface{}{"user_id": parameter["wuid"], "tags_id": val}
					save_arr = append(save_arr, marr)
				}
			}
		}
		res, err := DB().Table("merchant_wechat_user_addtags").Data(save_arr).Insert()
		if err != nil && len(save_arr) > 0 {
			results.Failed(context, "绑定分组失败", err)
		} else {
			results.Success(context, "绑定分组成功", res, nil)
		}
	}
}

//移除分组
func Removecate(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("merchant_wechat_user_addtags").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "移除失败", err)
	} else {
		results.Success(context, "移除成功！", res2, nil)
	}
	context.Abort()
	return
}

//编辑客户资料
func Upuserinfo(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	wuid := parameter["wuid"]
	delete(parameter, "wuid")
	res2, err := DB().Table("merchant_wechat_user").Where("wuid", wuid).Data(parameter).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		msg := "更新成功！"
		if res2 == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, res2, nil)
	}

}
