package packagecontent

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

//获取内容数据
func Getlist(context *gin.Context) {
	datatype := context.DefaultQuery("type", "goods")
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

	Fdb := DB().Table("m_shopping_goods_content").Where("accountID", user.Accountid).Where("type", datatype)
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
	getfield := "id,accountID,uid,cid,title,des,keywords,thumb,flag,price,original_price,browse,evaluate,weigh,status,createtime,useway,usetimes"
	if cid == 0 {
		_list, _err := Fdb.Fields(getfield).Page(pageNo).Limit(pageSize).Order("id desc").Get()
		list = _list
		err = _err
	} else {
		_cate_ids, _ := DB().Table("m_shopping_goods_cate").Where("pid", cid).Pluck("id")
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
		totalCount, _ = DB().Table("m_shopping_goods_content").Count()
		_pageSize := int64(pageSize)
		totalPage := totalCount / _pageSize
		for _, val := range list {
			catename, _ := DB().Table("m_shopping_goods_cate").Where("id", val["cid"]).Value("name")
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

//获取商品详情
func Getcontentmore(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	contentdata, _ := DB().Table("m_shopping_goods_content").Where("id", id).Fields("photo,des,content,instruction").First()
	contentpackage, _ := DB().Table("m_shopping_goods_content_package").Where("g_id", id).Get()
	contentdata["package"] = contentpackage
	results.Success(context, "获取详情", contentdata, nil)
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
	flagval, _ := utils.JsonMarshalNoSetEscapeHTML(parameter["flag"])
	parameter["flag"] = flagval
	photoval, _ := utils.JsonMarshalNoSetEscapeHTML(parameter["photo"])
	parameter["photo"] = photoval
	var packageitem interface{}
	if parameter["usetimes"] == "many" {
		packageitem = parameter["packageitem"]
		delete(parameter, "packageitem")
	}
	if f_id == 0 {
		parameter["accountID"] = user.Accountid
		parameter["createtime"] = time.Now().Unix()
		addId, err := DB().Table("m_shopping_goods_content").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			if addId != 0 {
				DB().Table("m_shopping_goods_content").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			if parameter["usetimes"] == "many" {
				add_packageitem(addId, packageitem)
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("m_shopping_goods_content").
			Data(parameter).
			Where("id", f_id).
			Update()
		if parameter["usetimes"] == "many" {
			add_packageitem(int64(f_id), packageitem)
		}
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			results.Success(context, "更新成功！", res, nil)
		}
	}
}

//批量提交-套餐
func add_packageitem(fid int64, packageitem interface{}) {
	save_arr := []map[string]interface{}{}
	for _, val := range packageitem.([]interface{}) {
		narr := val.(map[string]interface{})
		narr["g_id"] = fid
		save_arr = append(save_arr, narr)
	}
	DB().Table("m_shopping_goods_content_package").Where("g_id", fid).Delete()
	if len(save_arr) > 0 {
		DB().Table("m_shopping_goods_content_package").Data(save_arr).Insert()
	}
}

//删除内容-批量
func Delcontent_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("m_shopping_goods_content").WhereIn("id", ids.([]interface{})).Delete()
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
	res2, err := DB().Table("m_shopping_goods_content").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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

//删除套餐项目
func Delpackages_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	id := parameter["id"]
	res2, err := DB().Table("m_shopping_goods_content_package").Where("id", id).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}
