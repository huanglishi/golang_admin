package developers

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

//添加
func Add_cate_visiturl(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	f_id := parameter["id"].(float64)
	if f_id == 0 {
		addId, err := DB().Table("admin_dev_visiturl_cate").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("admin_dev_visiturl_cate").
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
func Get_cate_visiturl(context *gin.Context) {
	list, _ := DB().Table("admin_dev_visiturl_cate").Order("id asc").Get()
	for k, val := range list {
		getcount, _ := DB().Table("admin_dev_visiturl_content").Where("cid", val["id"]).Count()
		val["count"] = getcount
		if k == 0 {
			totalcount, _ := DB().Table("admin_dev_visiturl_content").Count()
			val["total"] = totalcount
		}
	}
	list_tree := GetTreeArray_only(list, 0)
	results.Success(context, "获取分类列表", list_tree, nil)
}

//获取父级数据
func Getparent_cate_visiturl(context *gin.Context) {
	list, _ := DB().Table("admin_dev_visiturl_cate").Where("status", 0).Fields("id,pid,name").Order("id asc").Get()
	rulenum := GetTreeArray(list, 0, "")
	list_text := GetTreeList_txt(rulenum, "name")
	results.Success(context, "获取父级数据", list_text, nil)
}

// 更新分类数据
func Up_cate_visiturl(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("admin_dev_visiturl_cate").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
func Del_cate_visiturl(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("admin_dev_visiturl_cate").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

//获取内容数据
func Getlist_visiturl(context *gin.Context) {
	_cid := context.DefaultQuery("cid", "0")
	titlekey := context.DefaultQuery("title", "")
	_pageNo := context.DefaultQuery("pageNo", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	cid, _ := strconv.Atoi(_cid)
	pageNo, _ := strconv.Atoi(_pageNo)
	pageSize, _ := strconv.Atoi(_pageSize)
	var list []gorose.Data
	var err error
	Fdb := DB().Table("admin_dev_visiturl_content")
	if titlekey != "" {
		Fdb = Fdb.Where("title", "like", "%"+titlekey+"%")
	}
	if cid == 0 {
		_list, _err := Fdb.Fields("id,cid,title,des,url,status,createtime").Page(pageNo).Limit(pageSize).Order("id desc").Get()
		list = _list
		err = _err
	} else {
		_cate_ids, _ := DB().Table("admin_dev_visiturl_cate").Where("pid", cid).Pluck("id")
		cate_ids := _cate_ids.([]interface{})
		cate_ids = append(cate_ids, cid)
		_list, _err := Fdb.WhereIn("cid", cate_ids).Fields("id,cid,title,des,url,status,createtime").Page(pageNo).Limit(pageSize).Order("id desc").Get()
		list = _list
		err = _err
	}
	if err != nil {
		results.Failed(context, "查找失败", err)
	} else {
		// 统计数据
		var totalCount int64
		totalCount, _ = DB().Table("admin_dev_visiturl_content").Count()
		_pageSize := int64(pageSize)
		totalPage := totalCount / _pageSize
		for _, val := range list {
			catename, _ := DB().Table("admin_dev_visiturl_cate").Where("id", val["cid"]).Value("name")
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

//添加内容
func Add_content_visiturl(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	f_id := parameter["id"].(float64)
	if f_id == 0 {
		parameter["createtime"] = time.Now().Unix()
		addId, err := DB().Table("admin_dev_visiturl_content").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("admin_dev_visiturl_content").
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
func Del_content_visiturl(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("admin_dev_visiturl_content").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 更新内容数据
func Up_content_visiturl(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("admin_dev_visiturl_content").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
