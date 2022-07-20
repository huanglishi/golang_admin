package business

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"math/rand"
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
		addId, err := DB().Table("merchant_user_cate").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("merchant_user_cate").
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
	list, _ := DB().Table("merchant_user_cate").Order("id asc").Get()
	for k, val := range list {
		getcount, _ := DB().Table("merchant_user").Where("cid", val["id"]).Where("type", 0).Count()
		val["count"] = getcount
		if k == 0 {
			totalcount, _ := DB().Table("merchant_user").Where("type", 0).Count()
			val["total"] = totalcount
		}
	}
	list_tree := GetTreeArray_only(list, 0)
	results.Success(context, "获取分类列表", list_tree, nil)
}

//获取父级数据
func Getparent_cate(context *gin.Context) {
	list, _ := DB().Table("merchant_user_cate").Where("status", 0).Fields("id,pid,name").Order("id asc").Get()
	rulenum := GetTreeArray(list, 0, "")
	list_text := GetTreeList_txt(rulenum, "name")
	//模块分类
	modulelist, _ := DB().Table("admin_dev_module_cate").Where("status", 0).Where("type", 0).Order("id asc").Get()
	results.Success(context, "获取父级数据", map[string]interface{}{"listcate": list_text, "modulelist": modulelist}, nil)
}

// 更新分类数据
func Up_cate(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("merchant_user_cate").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
	res2, err := DB().Table("merchant_user_cate").Where("id", parameter["id"]).Delete()
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
	namekey := context.DefaultQuery("name", "")
	_pageNo := context.DefaultQuery("pageNo", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	cid, _ := strconv.Atoi(_cid)
	pageNo, _ := strconv.Atoi(_pageNo)
	pageSize, _ := strconv.Atoi(_pageSize)
	var list []gorose.Data
	var err error
	Fdb := DB().Table("merchant_user").Where("type", 0)
	if namekey != "" {
		Fdb = Fdb.Where("name", "like", "%"+namekey+"%")
	}
	getfield := "id,cid,name,mobile,industry,city,cityval,address,status,username,avatar,validtime,remark,createtime"
	if cid == 0 {
		_list, _err := Fdb.Fields(getfield).Page(pageNo).Limit(pageSize).Order("id desc").Get()
		list = _list
		err = _err
	} else {
		_cate_ids, _ := DB().Table("merchant_user_cate").Where("pid", cid).Pluck("id")
		cate_ids := _cate_ids.([]interface{})
		cate_ids = append(cate_ids, cid)
		_list, _err := Fdb.WhereIn("cid", cate_ids).Fields(getfield).Page(pageNo).Limit(pageSize).Order("id desc").Get()
		list = _list
		err = _err
	}
	if err != nil {
		results.Failed(context, "查找失败", err)
	} else {
		// 统计数据
		var totalCount int64
		totalCount, _ = DB().Table("merchant_user").Count()
		_pageSize := int64(pageSize)
		totalPage := totalCount / _pageSize
		for _, val := range list {
			catename, _ := DB().Table("merchant_user_cate").Where("id", val["cid"]).Value("name")
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
		rnd := rand.New(rand.NewSource(6))
		salt := strconv.Itoa(rnd.Int())
		mdpass := parameter["password"].(string) + salt
		parameter["password"] = utils.Md5(mdpass)
		parameter["salt"] = salt
		parameter["createtime"] = time.Now().Unix()
		addId, err := DB().Table("merchant_user").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			//添加开户权限
			DB().Table("merchant_user_group_access").Data(map[string]interface{}{"uid": addId, "group_id": 1}).Insert()
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		if parameter["password"] != nil {
			rnd := rand.New(rand.NewSource(6))
			salt := strconv.Itoa(rnd.Int())
			mdpass := parameter["password"].(string) + salt
			parameter["password"] = utils.Md5(mdpass)
			parameter["salt"] = salt
		}
		parameter["updatetime"] = time.Now().Unix()
		res, err := DB().Table("merchant_user").
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
	res2, err := DB().Table("merchant_user").WhereIn("id", ids.([]interface{})).Delete()
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
	res2, err := DB().Table("merchant_user").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
