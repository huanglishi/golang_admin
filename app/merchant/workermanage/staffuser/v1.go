package staffuser

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

//获取内容数据
func Getlist(context *gin.Context) {
	_cid := context.DefaultQuery("cid", "0")
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

	Fdb := DB().Table("m_workermanage_staff_user").Where("accountID", user.Accountid)
	if createtime != "" {
		ids_arr := strings.Split(createtime, `,`)
		t1, _ := time.Parse("2006-01-02 15:04:05", ids_arr[0]+" 00:00:00")
		t2, _ := time.Parse("2006-01-02 15:04:05", ids_arr[1]+" 23:59:59")
		log.Println(t2.Unix())
		Fdb = Fdb.WhereBetween("createtime", []interface{}{t1.Unix(), t2.Unix()})
	}
	if keywords != "" {
		Fdb = Fdb.Where("name", "like", "%"+keywords+"%")
	}
	getfield := "id,accountID,uid,wuid,openid,name,sex,born,mobile,city,cityval,address,headimgurl,entryTime,position,jobnumber,professional,flag,weigh,status,account,createtime"
	if cid == 0 {
		_list, _err := Fdb.Fields(getfield).Page(pageNo).Limit(pageSize).Order("id desc").Get()
		list = _list
		err = _err
	} else {
		_cate_ids, _ := DB().Table("m_workermanage_staff_group").Where("pid", cid).Pluck("id")
		cate_ids := _cate_ids.([]interface{})
		cate_ids = append(cate_ids, cid)
		workerids, _ := DB().Table("m_workermanage_staff_group_bind").WhereIn("group_id", cate_ids).Pluck("worker_id")
		_workerids := workerids.([]interface{})
		if len(_workerids) > 0 {
			_list, _err := Fdb.WhereIn("id", _workerids).Fields(getfield).Page(pageNo).Limit(pageSize).Order("id desc").Get()
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
		totalCount, _ = DB().Table("m_workermanage_staff_user").Count()
		_pageSize := int64(pageSize)
		totalPage := totalCount / _pageSize
		for _, val := range list {
			catedata, _ := DB().Table("m_workermanage_staff_group_bind a").LeftJoin("m_workermanage_staff_group b on a.group_id = b.id").Where("a.worker_id", val["id"]).Fields("a.id,a.group_id,b.pid,b.name").Get()
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

//获取详情数据大的字段
func Getcontentmore_v1(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	contentdata, _ := DB().Table("m_workermanage_staff_user").Where("id", id).Fields("remark,skilled,content").First()
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
	//时间转化
	if parameter["entryTime"] != nil {
		loc, _ := time.LoadLocation("Local") //获取时区
		tmp, _ := time.ParseInLocation("2006-01-02", parameter["entryTime"].(string), loc)
		parameter["entryTime"] = tmp.Unix()
	}
	f_id := parameter["id"].(float64)
	if parameter["flag"] != nil {
		flagval, _ := utils.JsonMarshalNoSetEscapeHTML(parameter["flag"])
		parameter["flag"] = flagval
	}
	if f_id == 0 {
		parameter["accountID"] = user.Accountid
		parameter["createtime"] = time.Now().Unix()
		//处理密码
		rnd := rand.New(rand.NewSource(6))
		salt := strconv.Itoa(rnd.Int())
		mdpass := parameter["password"].(string) + salt
		parameter["password"] = utils.Md5(mdpass)
		parameter["salt"] = salt
		addId, err := DB().Table("m_workermanage_staff_user").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			if addId != 0 {
				DB().Table("m_workermanage_staff_user").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
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
		res, err := DB().Table("m_workermanage_staff_user").
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
	res2, err := DB().Table("m_workermanage_staff_user").WhereIn("id", ids.([]interface{})).Delete()
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
	res2, err := DB().Table("m_workermanage_staff_user").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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

// 绑定分组
func Setgroups_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter != nil {
		//批量提交
		save_arr := []map[string]interface{}{}
		for _, val := range parameter["groupids"].([]interface{}) {
			choosepack, _ := DB().Table("m_workermanage_staff_group_bind").Where("worker_id", parameter["uid"]).Where("group_id", val).Fields("id").Distinct().First()
			if choosepack == nil {
				catedata, _ := DB().Table("m_workermanage_staff_group").Where("pid", val).Fields("id").Distinct().First()
				if catedata == nil {
					marr := map[string]interface{}{"worker_id": parameter["uid"], "group_id": val}
					save_arr = append(save_arr, marr)
				}
			}
		}
		res, err := DB().Table("m_workermanage_staff_group_bind").Data(save_arr).Insert()
		if err != nil && len(save_arr) > 0 {
			results.Failed(context, "绑定分组失败", err)
		} else {
			results.Success(context, "绑定分组成功", res, nil)
		}
	}
}

//移除分组
func Removegroup_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("m_workermanage_staff_group_bind").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "移除失败", err)
	} else {
		results.Success(context, "移除成功！", res2, nil)
	}
	context.Abort()
	return
}
