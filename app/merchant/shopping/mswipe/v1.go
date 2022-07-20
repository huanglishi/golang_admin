package mswipe

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//获取内容数据
func Getlist(context *gin.Context) {
	title := context.DefaultQuery("title", "")
	keywords := context.DefaultQuery("keywords", "")
	createtime := context.DefaultQuery("createtime", "")
	_pageNo := context.DefaultQuery("pageNo", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(_pageNo)
	pageSize, _ := strconv.Atoi(_pageSize)
	//获取用户信息
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)

	Fdb := DB().Table("m_swipe_ads").Where("accountID", user.Accountid)
	if createtime != "" {
		ids_arr := strings.Split(createtime, `,`)
		t1, _ := time.Parse("2006-01-02 15:04:05", ids_arr[0]+" 00:00:00")
		t2, _ := time.Parse("2006-01-02 15:04:05", ids_arr[1]+" 23:59:59")
		log.Println(t2.Unix())
		Fdb = Fdb.WhereBetween("createtime", []interface{}{t1.Unix(), t2.Unix()})
	}
	if title != "" {
		Fdb = Fdb.Where("title", "like", "%"+title+"%")
	}
	if keywords != "" {
		Fdb = Fdb.Where("keywords", "like", "%"+keywords+"%")
	}
	getfield := "id,accountID,uid,type,title,image,link,linktype,valid_time,isskip,visits,likes,weigh,status,createtime"
	list, err := Fdb.Fields(getfield).Page(pageNo).Limit(pageSize).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "查找失败", err)
	} else {
		// 统计数据
		var totalCount int64
		totalCount, _ = DB().Table("m_swipe_ads").Count()
		_pageSize := int64(pageSize)
		totalPage := totalCount / _pageSize
		results.Success(context, "查找成功！", map[string]interface{}{
			"pageNo":     pageNo,
			"pageSize":   pageSize,
			"totalCount": totalCount,
			"totalPage":  totalPage,
			"data":       list}, nil)
	}
}

//获取文章详情
func Getcontentmore(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	contentdata, _ := DB().Table("m_swipe_ads").Where("id", id).Fields("content").First()
	results.Success(context, "获取文章详情", contentdata, nil)
}

//添加
func Pushcontent_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	f_id := parameter["id"].(float64)
	//数组转 字符串
	parameter["type"] = strings.Replace(strings.Trim(fmt.Sprint(parameter["type"]), "[]"), " ", ",", -1)
	//时间转化
	if parameter["valid_time"] != nil {
		loc, _ := time.LoadLocation("Local")
		valid_time, _ := time.ParseInLocation("2006-01-02", parameter["valid_time"].(string), loc)
		parameter["valid_time"] = valid_time.Unix()
	}
	if f_id == 0 {
		parameter["accountID"] = user.Accountid
		parameter["createtime"] = time.Now().Unix()
		addId, err := DB().Table("m_swipe_ads").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			if addId != 0 {
				DB().Table("m_swipe_ads").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("m_swipe_ads").
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

//删除内容-批量
func Delcontent_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("m_swipe_ads").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 更新用户状态
func Upstatus_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("m_swipe_ads").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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

//获取跳转链接素材
func Getmaterial(context *gin.Context) {
	_type := context.DefaultQuery("type", "article")
	listid := context.DefaultQuery("listid", "")
	searchkey := context.DefaultQuery("searchkey", "")
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	var list interface{}
	if _type == "article" {
		Fdb := DB().Table("m_marticle_content").Where("accountID", user.Accountid).Where("status", 0)
		if listid != "" {
			Fdb = Fdb.Where("id", "<", listid)
		}
		if searchkey != "" {
			Fdb = Fdb.Where("title", "like", "%"+searchkey+"%")
		}
		data, _ := Fdb.Fields("id,title").Limit(10).Order("id desc").Get()
		list = data
	} else {
		Fdb := DB().Table("m_shopping_goods_content").Where("accountID", user.Accountid).Where("status", 0)
		if listid != "" {
			Fdb = Fdb.Where("id", "<", listid)
		}
		if searchkey != "" {
			Fdb = Fdb.Where("title", "like", "%"+searchkey+"%")
		}
		data, _ := Fdb.Fields("id,title,type").Limit(10).Order("id desc").Get()
		list = data
	}
	results.Success(context, "获取跳转链接素材", list, nil)
}
