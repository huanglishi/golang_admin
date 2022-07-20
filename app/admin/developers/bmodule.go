package developers

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
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
		addId, err := DB().Table("admin_dev_module_cate").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("admin_dev_module_cate").
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
	list, _ := DB().Table("admin_dev_module_cate").Order("id asc").Get()
	for k, val := range list {
		getcount, _ := DB().Table("admin_dev_module_content").Where("cid", val["id"]).Count()
		val["count"] = getcount
		if k == 0 {
			totalcount, _ := DB().Table("admin_dev_module_content").Count()
			val["total"] = totalcount
		}
	}
	results.Success(context, "获取分类列表", list, nil)
}

//获取添加编辑-基础数据
func Getparent(context *gin.Context) {
	list, _ := DB().Table("admin_dev_module_cate").Order("id asc").Get()
	results.Success(context, "获取基础数据", list, nil)
}

// 更新分类数据
func Up_cate(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("admin_dev_module_cate").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
	res2, err := DB().Table("admin_dev_module_cate").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

//获取内容数据
func Getlist(context *gin.Context) {
	_cid := context.DefaultQuery("cid", "0")
	_pageNo := context.DefaultQuery("pageNo", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	cid, _ := strconv.Atoi(_cid)
	pageNo, _ := strconv.Atoi(_pageNo)
	pageSize, _ := strconv.Atoi(_pageSize)
	var list []gorose.Data
	var err error
	if cid == 0 {
		_list, _err := DB().Table("admin_dev_module_content").Fields("id,cid,title,des,image,status,weigh,createtime").Page(pageNo).Limit(pageSize).Order("weigh desc").Get()
		list = _list
		err = _err
	} else {
		_list, _err := DB().Table("admin_dev_module_content").Where("cid", cid).Fields("id,cid,title,des,image,status,weigh,createtime").Page(pageNo).Limit(pageSize).Order("weigh desc").Get()
		list = _list
		err = _err
	}
	if err != nil {
		results.Failed(context, "查找失败", err)
	} else {
		// 统计数据
		var totalCount int64
		totalCount, _ = DB().Table("admin_dev_module_content").Count()
		_pageSize := int64(pageSize)
		totalPage := totalCount / _pageSize
		for _, val := range list {
			catename, _ := DB().Table("admin_dev_module_cate").Where("id", val["cid"]).Value("name")
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
func Add_content(context *gin.Context) {
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
		addId, err := DB().Table("admin_dev_module_content").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			if addId != 0 {
				DB().Table("admin_dev_module_content").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			DB().Table("admin_dev_module_version").Data(map[string]interface{}{"mcid": addId, "name": "v1版本", "api": "v1", "weigh": 1, "des": parameter["title"].(string) + "基础版本"}).Insert()
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		parameter["updatetime"] = time.Now().Unix()
		res, err := DB().Table("admin_dev_module_content").
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
func Del_content(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("admin_dev_module_content").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 更新内容数据
func Up_content(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("admin_dev_module_content").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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

//添加版本
func Addversion(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	f_id := parameter["id"].(float64)
	delete(parameter, "editable")
	if f_id == 0 {
		delete(parameter, "isNew")
		addId, err := DB().Table("admin_dev_module_version").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			if addId != 0 {
				DB().Table("admin_dev_module_version").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		delete(parameter, "_originalData")
		res, err := DB().Table("admin_dev_module_version").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else if res == 1 {
			results.Success(context, "更新成功！", f_id, nil)
		} else {
			results.Success(context, "没有需要更新内容！", f_id, nil)
		}
	}
}

//获取版本数据
func Getversion(context *gin.Context) {
	_mcid := context.DefaultQuery("mcid", "0")
	mcid, _ := strconv.Atoi(_mcid)
	list, _ := DB().Table("admin_dev_module_version").Where("mcid", mcid).Order("id asc").Get()
	for _, val := range list {
		val["editable"] = false
	}
	results.Success(context, "获取版本数据", list, nil)
}

//删除版本数据
func Delversion(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("admin_dev_module_version").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 更新版本数据
func Upversion(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("admin_dev_module_version").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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

//获取版本描述
func Getversiondes(context *gin.Context) {
	_id := context.DefaultQuery("id", "0")
	id, _ := strconv.Atoi(_id)
	des, _ := DB().Table("admin_dev_module_version").Where("id", id).Value("des")
	results.Success(context, "获取版本描述", des, nil)
}

// 更新版本描述
func Upversiondes(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res, err := DB().Table("admin_dev_module_version").Where("id", parameter["id"]).Data(map[string]interface{}{"des": parameter["des"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		msg := "更新成功！"
		if res == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, res, nil)
	}
}

//获取菜单权限
func Getversionmenu(context *gin.Context) {
	//获取默认值
	id := context.DefaultQuery("id", "0")
	datafield, _ := DB().Table("admin_dev_module_version").Where("id", id).Fields("rules,checkedKeys").First()
	menu, _ := DB().Table("merchant_client_menu").Where("isbase", 1).Where("status", 0).Fields("id,pid,weigh,title").Order("weigh asc").Get()
	for _, mv := range menu {
		mv["scopedSlots"] = map[string]interface{}{"title": "title"}
	}
	menu_tree := GetTreeArray_only(menu, 0)
	var parameter = map[string]interface{}{"version": datafield, "menu": menu_tree, "menuarr": menu}
	results.Success(context, "获取基础数据", parameter, nil)
}

// 更新版本权限菜单
func Upversionmenu(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res, err := DB().Table("admin_dev_module_version").Where("id", parameter["id"]).Data(map[string]interface{}{"rules": parameter["rules"], "checkedKeys": parameter["checkedKeys"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		msg := "更新成功！"
		if res == 0 {
			msg = "暂无数据更新"
		} else {
			//更新菜单版本
			nemu_ids := strings.Split(parameter["rules"].(string), `,`)
			api_version, err := DB().Table("admin_dev_module_version").Where("id", parameter["id"]).Value("api")
			if err == nil {
				for _, menu_id := range nemu_ids {
					DB().Table("merchant_client_menu").Where("id", menu_id).Data(map[string]interface{}{"api_version": api_version}).Update()
				}
			}
		}
		results.Success(context, msg, res, nil)
	}
}

//获取前台web访问地址
func Getvisiturl(context *gin.Context) {
	//获取默认值
	id := context.DefaultQuery("id", "0")
	visiturl_ids, _ := DB().Table("admin_dev_module_version").Where("id", id).Value("visiturl_ids")
	orlist, _ := DB().Table("admin_dev_visiturl_cate").Fields("id,pid,name").Order("id asc").Get()
	for _, mv := range orlist {
		clist, _ := DB().Table("admin_dev_visiturl_content").Where("cid", mv["id"]).Where("status", 0).Fields("id,cid,title,des,url").Order("id asc").Get()
		mv["list"] = clist
	}
	orlist_tree := GetTreeArray_only(orlist, 0)
	var parameter = map[string]interface{}{"orlist": orlist_tree, "visiturlids": visiturl_ids}
	results.Success(context, "获取地址数据", parameter, nil)
}

//更新前台web访问地址数据
func Upvisiturl(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res, err := DB().Table("admin_dev_module_version").Where("id", parameter["id"]).Data(map[string]interface{}{"visiturl_ids": parameter["visiturl"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		results.Success(context, "更新成功！", res, nil)
	}
}
