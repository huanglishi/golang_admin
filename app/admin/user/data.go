package user

import (
	"basegin/app/model"
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

func DB() gorose.IOrm {
	return model.DB.NewOrm()
}

/*
  1.获取用户信息
  2.获取用户的
*/
func GetInfo(context *gin.Context) {
	getuser, ok := context.Get("user") //取值 实现了跨中间件取值
	if !ok {
		results.Failed(context, "用户id不存在！", ok)
		return
	}
	user := getuser.(*utils.UserClaims)
	res, err := DB().Table("admin_user").Fields("id,name,avatar,telephone,email,lastLoginIp,lastLoginTime,status,valid_time").Where("id", user.ID).First()
	// res, err := DB().Table("admin_user").Fields("id, name").Limit(10).Offset(0).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "查找用户数据！", err)
		return
	}
	//添加权限
	group_ids, _ := DB().Table("admin_auth_group_access").Where("uid", user.ID).Pluck("group_id")
	roles_ids, _ := DB().Table("admin_auth_group").WhereIn("id", group_ids.([]interface{})).Pluck("rules")
	//先处理权限数据
	var ids_arr_int []interface{}
	if ok, _ = utils.Contain("*", roles_ids); !ok {
		var rule_ids_arr []string
		for _, v := range roles_ids.([]interface{}) {
			if v != nil || v != "" {
				ids_arr := strings.Split(v.(string), `,`)
				rule_ids_arr = append(rule_ids_arr, ids_arr...)
			}
		}
		ids_buy, _ := json.Marshal(&rule_ids_arr)
		json.Unmarshal([]byte(ids_buy), &ids_arr_int)
	} else {
		ids_arr_int = make([]interface{}, 0)
	}
	// b, _ := json.Marshal(&res)
	// var m map[string]interface{}
	// _ = json.Unmarshal(b, &m)
	//2.获取权限数据
	var _rule []gorose.Data
	if ok, _ = utils.Contain("*", roles_ids); ok {
		rule, _ := DB().Table("admin_auth_rule").Where("type", "menu").Fields("id,type,pid,icon,weigh,title,status,component").AddFields("name as permissionId").Order("weigh asc").Get()
		_rule = rule
	} else {
		rule, _ := DB().Table("admin_auth_rule").Where("type", "menu").WhereIn("id", ids_arr_int).Fields("id,type,pid,icon,weigh,title,status,component").AddFields("name as permissionId").Order("weigh asc").Get()
		_rule = rule
	}
	for _, v := range _rule {
		v["roleId"] = "admin"
		v["actionList"] = nil
		var filedata []gorose.Data
		if ok, _ = utils.Contain("*", roles_ids); ok {
			_filedata, _ := DB().Table("admin_auth_rule").Where("type", "file").Where("pid", v["id"]).Fields("id,title").AddFields("name as action").Order("weigh asc").Get()
			filedata = _filedata
		} else {
			_filedata, _ := DB().Table("admin_auth_rule").Where("type", "file").Where("pid", v["id"]).WhereIn("id", ids_arr_int).Fields("id,title").AddFields("name as action").Order("weigh asc").Get()
			filedata = _filedata
		}
		v["actionEntitySet"] = filedata
	}
	role := map[string]interface{}{"permissions": _rule}
	//3.菜单数据
	var _menu []gorose.Data
	if ok, _ = utils.Contain("*", roles_ids); ok {
		menu, _ := DB().Table("admin_auth_rule").Where("type", "menu").Where("status", 0).Fields("id,pid,weigh,title,name,component,redirect,target,icon,keepAlive,hideHeader,hiddenHeaderContent,hideChildrenInMenu,hidden").Order("weigh asc").Get()
		_menu = menu
	} else {
		menu, _ := DB().Table("admin_auth_rule").Where("type", "menu").Where("status", 0).WhereIn("id", ids_arr_int).Fields("id,pid,weigh,title,name,component,redirect,target,icon,keepAlive,hideHeader,hiddenHeaderContent,hideChildrenInMenu,hidden").Order("weigh asc").Get()
		_menu = menu
	}
	menulist := GetTreeArray(_menu, 0)
	//返回登录信息
	context.JSON(200, gin.H{
		"code":    0,
		"msg":     "查找用户数据成功！",
		"info":    res,
		"role":    role,
		"menu":    menulist,
		"gruopid": group_ids,
		"time":    time.Now().Unix(),
	})
}

// 更新数据
func Updata(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("admin_user").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
	if err != nil {
		context.JSON(200, gin.H{
			"code": 1,
			"msg":  "更新失败！",
			"data": err,
			"time": time.Now().Unix(),
		})
	} else {
		msg := "更新成功！"
		if res2 == 0 {
			msg = "暂无数据更新"
		}
		context.JSON(200, gin.H{
			"code": 0,
			"msg":  msg,
			"data": res2,
			"time": time.Now().Unix(),
		})
	}

}

// 添加用户
func AddParam(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	parameter["createtime"] = time.Now().Unix()
	if parameter["valid_time"] != nil {
		valid_time, _ := time.Parse("01/02/2006", parameter["valid_time"].(string))
		parameter["valid_time"] = valid_time
	} else {
		parameter["valid_time"] = 0
	}
	f_id := parameter["id"].(float64)
	groupids := parameter["groupids"].([]interface{})
	delete(parameter, "groupids")
	if f_id == 0 {
		rnd := rand.New(rand.NewSource(6))
		salt := strconv.Itoa(rnd.Int())
		mdpass := parameter["password"].(string) + salt
		parameter["password"] = utils.Md5(mdpass)
		parameter["salt"] = salt
		addId, err := DB().Table("admin_user").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			group_arr := []map[string]interface{}{}
			for _, val := range groupids {
				marr := map[string]interface{}{"uid": addId, "group_id": val}
				group_arr = append(group_arr, marr)
			}
			DB().Table("admin_auth_group_access").Where("uid", f_id).Delete()
			DB().Table("admin_auth_group_access").Data(group_arr).Insert()
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
		res, err := DB().Table("admin_user").
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
			DB().Table("admin_auth_group_access").Where("uid", f_id).Delete()
			DB().Table("admin_auth_group_access").Data(group_arr).Insert()
			results.Success(context, "更新成功！", res, group_arr)
		}
	}
	context.Abort()
	return
}

// 获取用列表
func QueryParam(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	group_ids, _ := DB().Table("admin_auth_group_access").Where("uid", user.ID).Pluck("group_id")                      //获取用户分组
	data_access_ids, _ := DB().Table("admin_auth_group").WhereIn("id", group_ids.([]interface{})).Pluck("data_access") //获取用户数据权限
	_pageNo := context.DefaultQuery("pageNo", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(_pageNo)
	pageSize, _ := strconv.Atoi(_pageSize)
	//权限判断
	var res []gorose.Data
	var err error
	var access_all int64 = 2
	var access_myandchrid int64 = 1
	if utils.In_array(access_all, data_access_ids) { //全部
		_res, _err := DB().Table("admin_user").Fields("id, name,username,telephone,lastLoginIp,lastLoginTime,status,valid_time").Page(pageNo).Limit(pageSize).Order("id desc").Get()
		res = _res
		err = _err
	} else if utils.In_array(access_myandchrid, data_access_ids) { //自己和子集
		_res, _err := DB().Table("admin_user").Fields("id, name,username,telephone,lastLoginIp,lastLoginTime,status,valid_time").Page(pageNo).Limit(pageSize).Order("id desc").Get()
		res = _res
		err = _err
	} else { //仅自己
		_res, _err := DB().Table("admin_user").Where("uid", user.ID).Fields("id, name,username,telephone,lastLoginIp,lastLoginTime,status,valid_time").Page(pageNo).Limit(pageSize).Order("id desc").Get()
		res = _res
		err = _err
	}
	if err != nil {
		results.Failed(context, "查找失败", err)
	} else {
		// 统计数据
		var totalCount int64
		totalCount, _ = DB().Table("admin_user").Count()
		_pageSize := int64(pageSize)
		totalPage := totalCount / _pageSize
		for _, val := range res {
			groupids, _ := DB().Table("admin_auth_group_access").Where("uid", val["id"]).Pluck("group_id")
			val["groupids"] = groupids
			groupname, _ := DB().Table("admin_auth_group").WhereIn("id", groupids.([]interface{})).Fields("id,name").Get()
			val["groupname"] = groupname
		}
		results.Success(context, "查找成功！", map[string]interface{}{
			"pageNo":     pageNo,
			"pageSize":   pageSize,
			"totalCount": totalCount,
			"totalPage":  totalPage,
			"data":       res}, nil)
	}
}

//删除操作
func DelParam(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("admin_user").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return

}

//修改密码
func Changepassword(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//账号信息
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	userdata, err := DB().Table("admin_user").Where("id", user.ID).Fields("password,salt").First()
	if err != nil {
		results.Failed(context, "修改密码失败", err)
	} else {
		pass := utils.Md5(parameter["originalpassword"].(string) + userdata["salt"].(string))
		if userdata["password"] != pass {
			results.Failed(context, "原来密码输入错误！", err)
		} else {
			newpass := utils.Md5(parameter["newpassword"].(string) + userdata["salt"].(string))
			res, err := DB().Table("admin_user").
				Data(map[string]interface{}{"password": newpass}).
				Where("id", user.ID).
				Update()
			if err != nil {
				results.Failed(context, "修改密码失败", err)
			} else {
				results.Success(context, "修改密码成功！", res, nil)
			}
		}
	}
	context.Abort()
	return
}

//tool-获取树状数组
func GetTreeArray(num []gorose.Data, pid int64) []gorose.Data {
	childs := ToolFar(num, pid) //获取pid下的所有数据
	var chridnum []gorose.Data
	if childs != nil {
		for _, v := range childs {
			v["children"] = GetTreeArray(num, v["id"].(int64))
			chridnum = append(chridnum, v)
		}
	}
	return chridnum
}

//base_tool-获取pid下所有数组
func ToolFar(data []gorose.Data, pid int64) []gorose.Data {
	var mapString []gorose.Data
	for _, v := range data {
		if v["pid"].(int64) == pid {
			mapString = append(mapString, v)
		}
	}
	return mapString
}
