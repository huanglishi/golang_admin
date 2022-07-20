package auths

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

//获取内容数据
func Get_account(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	group_ids, _ := DB().Table("merchant_user_group_access").Where("uid", user.ID).Pluck("group_id")                      //获取用户分组
	data_access_ids, _ := DB().Table("merchant_user_group").WhereIn("id", group_ids.([]interface{})).Pluck("data_access") //获取用户数据权限
	_pageNo := context.DefaultQuery("pageNo", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(_pageNo)
	pageSize, _ := strconv.Atoi(_pageSize)
	//权限判断
	var res []gorose.Data
	var err error
	var totalCount int64
	var access_all int64 = 2
	var access_myandchrid int64 = 1
	FDB := DB().Table("merchant_user").Fields("id,pid, name,username,mobile,lastLoginIp,lastLoginTime,status,validtime").Where("accountID", user.Accountid)
	FDB_where := DB().Table("merchant_user").Where("accountID", user.Accountid)
	//判断数据权限
	if utils.In_array(access_all, data_access_ids) { //全部
		_res, _err := FDB.Page(pageNo).Limit(pageSize).Order("id desc").Get()
		res = _res
		err = _err
		totalCount, _ = FDB_where.Count()
	} else if utils.In_array(access_myandchrid, data_access_ids) { //自己和子集
		alldata, _ := FDB_where.Fields("id,pid").Order("id desc").Get()
		child_ids := Getallchildren(alldata, user.ID) //获取pid下的所有数据
		child_ids_arr := gettree_to_arr(child_ids)
		var chriddis []interface{}
		for _, grouparr := range child_ids_arr {
			chriddis = append(chriddis, grouparr["id"])
		}
		if len(chriddis) == 0 {
			chriddis = append(chriddis, 0)
		}
		_res, _err := FDB.WhereIn("id", chriddis).Page(pageNo).Limit(pageSize).Order("id desc").Get()
		res = _res
		err = _err
		totalCount, _ = FDB_where.WhereIn("id", chriddis).Count()
	} else { //仅自己
		_res, _err := FDB.Where("uid", user.ID).Page(pageNo).Limit(pageSize).Order("id desc").Get()
		res = _res
		err = _err
		totalCount, _ = FDB_where.Where("uid", user.ID).Count()
	}
	if err != nil {
		results.Failed(context, "查找失败", err)
	} else {
		// 统计数据
		_pageSize := int64(pageSize)
		totalPage := totalCount / _pageSize
		for _, val := range res {
			groupids, _ := DB().Table("merchant_user_group_access").Where("uid", val["id"]).Pluck("group_id")
			val["groupids"] = groupids
			if groupids != nil {
				groupname, _ := DB().Table("merchant_user_group").WhereIn("id", groupids.([]interface{})).Fields("id,name").Get()
				val["groupname"] = groupname
			} else {
				val["groupname"] = nil
			}
		}
		results.Success(context, "查找成功！", map[string]interface{}{
			"pageNo":     pageNo,
			"pageSize":   pageSize,
			"totalCount": totalCount,
			"totalPage":  totalPage,
			"data":       res}, nil)
	}
}

//添加内容
func Add_account(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)

	if parameter["validtime"] != nil {
		validtime, _ := time.Parse("01/02/2006", parameter["validtime"].(string))
		parameter["validtime"] = validtime
	} else {
		parameter["validtime"] = 0
	}
	f_id := parameter["id"].(float64)
	groupids := parameter["groupids"].([]interface{})
	delete(parameter, "groupids")
	if f_id == 0 {
		parameter["createtime"] = time.Now().Unix()
		parameter["uid"] = user.ID
		parameter["pid"] = user.ID
		parameter["type"] = 1
		parameter["accountID"] = user.Accountid
		rnd := rand.New(rand.NewSource(6))
		salt := strconv.Itoa(rnd.Int())
		mdpass := parameter["password"].(string) + salt
		parameter["password"] = utils.Md5(mdpass)
		parameter["salt"] = salt
		addId, err := DB().Table("merchant_user").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			group_arr := []map[string]interface{}{}
			for _, val := range groupids {
				marr := map[string]interface{}{"uid": addId, "group_id": val}
				group_arr = append(group_arr, marr)
			}
			DB().Table("merchant_user_group_access").Where("uid", f_id).Delete()
			DB().Table("merchant_user_group_access").Data(group_arr).Insert()
			results.Success(context, "添加成功！", nil, nil)
		}
	} else {
		if parameter["password"] != nil {
			rnd := rand.New(rand.NewSource(6))
			salt := strconv.Itoa(rnd.Int())
			mdpass := parameter["password"].(string) + salt
			parameter["password"] = utils.Md5(mdpass)
			parameter["salt"] = salt
		}
		res, err := DB().Table("merchant_user").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			group_arr := []map[string]interface{}{}
			for _, val := range groupids {
				marr := map[string]interface{}{"uid": f_id, "group_id": val}
				group_arr = append(group_arr, marr)
			}
			DB().Table("merchant_user_group_access").Where("uid", f_id).Delete()
			DB().Table("merchant_user_group_access").Data(group_arr).Insert()
			results.Success(context, "更新成功！", res, group_arr)
		}
	}
	context.Abort()
	return
}

//删除内容-批量
func Del_account(context *gin.Context) {
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
func Up_account(context *gin.Context) {
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
