package leavemessage

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

//获取获取网站列表
func GetSitelist_v1(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("m_website_site").Where("accountID", user.Accountid).Fields("id,name,status").Order("id asc").Get()
	for k, val := range list {
		getcount, _ := DB().Table("m_website_site_leavemessage").Where("site_id", val["id"]).Count()
		val["count"] = getcount
		if k == 0 {
			totalcount, _ := DB().Table("m_website_site_leavemessage").Where("accountID", user.Accountid).Count()
			val["total"] = totalcount
		}
	}
	results.Success(context, "获取网站列表", list, nil)
}

//获取内容数据
func Getlist_v1(context *gin.Context) {
	_cid := context.DefaultQuery("cid", "0")
	title := context.DefaultQuery("title", "")
	keywords := context.DefaultQuery("keywords", "")
	createtime := context.DefaultQuery("createtime", "")
	_pageNo := context.DefaultQuery("pageNo", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	cid, _ := strconv.Atoi(_cid)
	pageNo, _ := strconv.Atoi(_pageNo)
	pageSize, _ := strconv.Atoi(_pageSize)
	var list []gorose.Data
	var err error
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)

	Fdb := DB().Table("m_website_site_leavemessage").Where("accountID", user.Accountid)
	if createtime != "" {
		ids_arr := strings.Split(createtime, `,`)
		t1, _ := time.Parse("2006-01-02 15:04:05", ids_arr[0]+" 00:00:00")
		t2, _ := time.Parse("2006-01-02 15:04:05", ids_arr[1]+" 23:59:59")
		log.Println(t2.Unix())
		Fdb = Fdb.WhereBetween("createtime", []interface{}{t1.Unix(), t2.Unix()})
	}
	if title != "" {
		Fdb = Fdb.Where("name", "like", "%"+title+"%")
	}
	if keywords != "" {
		Fdb = Fdb.Where("content", "like", "%"+keywords+"%")
	}
	getfield := "id,name,site_id,mobile,status,createtime"
	if cid == 0 {
		_list, _err := Fdb.Fields(getfield).Page(pageNo).Limit(pageSize).Order("id desc").Get()
		list = _list
		err = _err
	} else {
		_list, _err := Fdb.Where("site_id", cid).Fields(getfield).Page(pageNo).Limit(pageSize).Order("id desc").Get()
		list = _list
		err = _err
	}
	if err != nil {
		results.Failed(context, "查找失败", err)
	} else {
		// 统计数据
		var totalCount int64
		totalCount, _ = DB().Table("m_website_site_leavemessage").Count()
		_pageSize := int64(pageSize)
		totalPage := totalCount / _pageSize
		for _, val := range list {
			catename, _ := DB().Table("m_website_site").Where("id", val["site_id"]).Value("name")
			val["catename"] = catename
		}
		results.Success(context, "查找成功！", map[string]interface{}{
			"pageNo":     pageNo,
			"pageSize":   pageSize,
			"totalCount": totalCount,
			"totalPage":  totalPage,
			"data":       list}, nil)
	}
}

//获取文章详情
func Getcontentmore_v1(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	contentdata, _ := DB().Table("m_website_site_leavemessage").Where("id", id).Fields("content").First()
	DB().Table("m_website_site_leavemessage").Where("id", id).Data(map[string]interface{}{"status": 1}).Update()
	results.Success(context, "获取文章详情", contentdata, nil)
}

//删除内容-批量
func Delcontent_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("m_website_site_leavemessage").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 更新状态
func Upstatus_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("m_website_site_leavemessage").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
