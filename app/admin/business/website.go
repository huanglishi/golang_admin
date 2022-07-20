package business

import (
	"basegin/utils/results"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

//获取网站列表
func Websitelist(context *gin.Context) {
	_cid := context.DefaultQuery("cid", "0")
	namekey := context.DefaultQuery("name", "")
	_pageNo := context.DefaultQuery("pageNo", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	cid, _ := strconv.Atoi(_cid)
	pageNo, _ := strconv.Atoi(_pageNo)
	pageSize, _ := strconv.Atoi(_pageSize)
	var list []gorose.Data
	var err error
	Fdb := DB().Table("m_website_site")
	if namekey != "" {
		Fdb = Fdb.Where("name", "like", "%"+namekey+"%")
	}
	getfield := "id,accountID,name,domain,thumb,des,valid_time,status,step,company,mobile,createtime"
	if cid == 0 {
		_list, _err := Fdb.Fields(getfield).Page(pageNo).Limit(pageSize).Order("id desc").Get()
		list = _list
		err = _err
	} else {
		_cate_ids, _ := DB().Table("merchant_user_cate").Where("pid", cid).Pluck("id")
		cate_ids := _cate_ids.([]interface{})
		cate_ids = append(cate_ids, cid)
		user_ids, _ := DB().Table("merchant_user").WhereIn("cid", cate_ids).Pluck("id")
		_list, _err := Fdb.WhereIn("accountID", user_ids.([]interface{})).Fields(getfield).Page(pageNo).Limit(pageSize).Order("id desc").Get()
		list = _list
		err = _err
	}
	if err != nil {
		results.Failed(context, "查找失败", err)
	} else {
		// 统计数据
		var totalCount int64
		totalCount, _ = DB().Table("m_website_site").Count()
		_pageSize := int64(pageSize)
		totalPage := totalCount / _pageSize
		for _, val := range list {
			//1获取用户
			muser, _ := DB().Table("merchant_user").Where("id", val["accountID"]).Fields("cid,mobile,name").First()
			catename, _ := DB().Table("merchant_user_cate").Where("id", muser["cid"]).Value("name")
			val["catename"] = catename
			val["m_mobile"] = muser["mobile"]
			val["m_name"] = muser["name"]
		}
		results.Success(context, "查找成功！", map[string]interface{}{
			"pageNo":     pageNo,
			"pageSize":   pageSize,
			"totalCount": totalCount,
			"totalPage":  totalPage,
			"data":       list}, nil)
	}
}

// 审批
func Up_audit(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("m_website_site").WhereIn("id", ids_arr).Data(map[string]interface{}{"step": parameter["step"]}).Update()
	if err != nil {
		results.Failed(context, "审批失败！", err)
	} else {
		msg := "审批成功！"
		if res2 == 0 {
			msg = "暂无数据审批"
		}
		results.Success(context, msg, res2, nil)
	}
}

//获取网站指定字段数据
func Getsitefield_v1(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	fields := context.DefaultQuery("fields", "")
	data, _ := DB().Table("m_website_site").Where("id", id).Fields(fields).First()
	results.Success(context, "获取网站指定字段数据", data, nil)
}

//更新网站配置
func Upsite_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res, err := DB().Table("m_website_site").
		Data(parameter).
		Where("id", parameter["id"]).
		Update()
	if err != nil {
		results.Failed(context, "更新失败", err)
	} else {
		results.Success(context, "更新成功！", res, nil)
	}
}

//获取所有网站列表
func Sitelist_v1(context *gin.Context) {
	namekey := context.DefaultQuery("name", "")
	id := context.DefaultQuery("id", "0")
	Fdb := DB().Table("m_website_site")
	if namekey != "" {
		Fdb = Fdb.Where("name", "like", "%"+namekey+"%")
	}
	if id != "0" {
		Fdb = Fdb.Where("id", "!=", id)
	}
	list, _ := Fdb.Fields("id,name,accountID").Order("id asc").Get()
	for _, val := range list {
		//1获取用户
		muser, _ := DB().Table("merchant_user").Where("id", val["accountID"]).Fields("cid,name").First()
		// catename, _ := DB().Table("merchant_user_cate").Where("id", muser["cid"]).Value("name")
		// val["catename"] = catename
		val["m_name"] = muser["name"]
	}
	results.Success(context, "获取网站列表", list, nil)
}

//清除原来网站数据
func Clearsite_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	id := parameter["id"]
	// 开始事务
	db := DB()
	db.Begin()
	_, err_m := db.Table("m_website_module").Where("site_id", id).Delete()
	if err_m != nil {
		// 回滚事务
		db.Rollback()
		results.Failed(context, "清除失败", err_m)
	} else {
		_, err_cate := db.Table("m_website_article_cate").Where("site_id", id).Delete()
		if err_cate != nil {
			// 回滚事务
			db.Rollback()
			results.Failed(context, "清除失败", err_cate)
		} else {
			_, err_ac := db.Table("m_website_article_content").Where("site_id", id).Delete()
			if err_ac != nil {
				// 回滚事务
				db.Rollback()
				results.Failed(context, "清除失败", err_ac)
			} else {
				// 提交事务
				db.Commit()
				results.Success(context, "清除网站成功", id, nil)
			}
		}
	}
}

//克隆指定网站数据
func Clonesite_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	id := parameter["id"].(float64) //当前网站id
	site_id := parameter["site_id"].(float64)
	accountID := parameter["accountID"].(float64)
	//[]interface{}=[]string
	clonedata := parameter["clonedata"].([]interface{})
	array := make([]string, len(clonedata))
	for i, v := range clonedata {
		array[i] = v.(string)
	}
	if IsContain(array, "module") { //克隆模块
		//读取数据
		mlist, _ := DB().Table("m_website_module").Where("site_id", site_id).Fields("accountID,uid,site_id,name,model_name,template_list,template_show,singlepage,fcatdir").Get()
		for _, val := range mlist {
			val["site_id"] = id
			val["accountID"] = accountID
		}
		_, err_m := DB().Table("m_website_module").Data(mlist).Insert()
		if err_m != nil {
			results.Failed(context, "克隆模块失败", err_m)
		} else {
			results.Success(context, "克隆成功", nil, nil)
		}
	}

}
