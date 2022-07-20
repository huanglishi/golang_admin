package user

import (
	"basegin/app/model"
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

func DB() gorose.IOrm {
	model.DB.GetExecuteDB().SetConnMaxLifetime(time.Duration(8*3600) * time.Second) //失效时间
	return model.DB.NewOrm()
}

/*
  1.获取用户信息
  2.获取用户的
*/
func GetInfo(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	//获取用户信息
	userinfo, err := DB().Table("merchant_user").Fields("id,pid,name,avatar,mobile,lastLoginIp,lastLoginTime,status,remark,validtime").Where("id", user.ID).First()
	if err != nil {
		results.Failed(context, "查找用户数据！", err)
		return
	}
	//添加权限
	group_ids, _ := DB().Table("merchant_user_group_access").Where("uid", user.ID).Pluck("group_id")
	roles_ids, _ := DB().Table("merchant_user_group").WhereIn("id", group_ids.([]interface{})).Pluck("rules")
	//定义权限数组
	var role map[string]interface{}
	var ids_arr_int []interface{}
	if roles_ids != nil {
		//先处理权限数据
		if !utils.In_array("*", roles_ids) {
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
		//购买的权限
		choosepack, _ := DB().Table("admin_business_choosepack").Where("validtime", ">", time.Now().Unix()).Where("ispay", 1).Where("accountID", user.ID).Pluck("vid")
		choosepack_or, _ := DB().Table("admin_business_choosepack").Where("validtime", 0).Where("ispay", 1).Where("accountID", user.ID).Pluck("vid")
		choosepack_arr := choosepack.([]interface{}) //转换
		for _, val := range choosepack_or.([]interface{}) {
			choosepack_arr = append(choosepack_arr, val)
		}
		roles_version, _ := DB().Table("admin_dev_module_version").WhereIn("id", choosepack_arr).Pluck("rules")
		var roles_version_int []interface{}
		if roles_version != nil {
			//先处理权限数据
			if !utils.In_array("*", roles_version) {
				var roles_version_arr []string
				for _, v := range roles_version.([]interface{}) {
					if v != nil || v != "" {
						version_ids_arr := strings.Split(v.(string), `,`)
						roles_version_arr = append(roles_version_arr, version_ids_arr...)
					}
				}
				ids_buy, _ := json.Marshal(&roles_version_arr)
				json.Unmarshal([]byte(ids_buy), &roles_version_int)
			} else {
				roles_version_int = make([]interface{}, 0)
			}
		}
		//2.获取权限数据
		var _rule []gorose.Data
		if utils.In_array("*", roles_ids) {
			if roles_version != nil {
				rule, _ := DB().Table("merchant_client_menu").Where("type", "menu").Where("isbase", 0).OrWhereIn("id", roles_version_int).Fields("id,type,pid,icon,weigh,title,status,component").AddFields("name as permissionId").Order("weigh asc").Get()
				_rule = rule
			} else {
				rule, _ := DB().Table("merchant_client_menu").Where("type", "menu").Where("isbase", 0).Fields("id,type,pid,icon,weigh,title,status,component").AddFields("name as permissionId").Order("weigh asc").Get()
				_rule = rule
			}
		} else {
			rule, _ := DB().Table("merchant_client_menu").Where("type", "menu").WhereIn("id", ids_arr_int).Fields("id,type,pid,icon,weigh,title,status,component").AddFields("name as permissionId").Order("weigh asc").Get()
			_rule = rule
		}
		for _, v := range _rule {
			v["actionList"] = nil
			var filedata []gorose.Data
			if utils.In_array("*", roles_ids) {
				if roles_version != nil {
					_filedata, _ := DB().Table("merchant_client_menu").Where("type", "file").Where("pid", v["id"]).Where("isbase", 0).OrWhereIn("id", roles_version_int).Fields("id,title").AddFields("name as action").Order("weigh asc").Get()
					filedata = _filedata
				} else {
					_filedata, _ := DB().Table("merchant_client_menu").Where("type", "file").Where("pid", v["id"]).Where("isbase", 0).Fields("id,title").AddFields("name as action").Order("weigh asc").Get()
					filedata = _filedata
				}
			} else {
				_filedata, _ := DB().Table("merchant_client_menu").Where("type", "file").Where("pid", v["id"]).WhereIn("id", ids_arr_int).Fields("id,title").AddFields("name as action").Order("weigh asc").Get()
				filedata = _filedata
			}
			v["actionEntitySet"] = filedata
		}
		//组装权限
		role = map[string]interface{}{"permissions": _rule}
		//3.菜单数据
		var menu []gorose.Data
		Menu_Com := DB().Table("merchant_client_menu").Where("type", "menu").Where("status", 0).Fields("id,pid,weigh,title,name,component,redirect,target,icon,keepAlive,hideHeader,hiddenHeaderContent,hideChildrenInMenu,hidden,api_version")
		if utils.In_array("*", roles_ids) {
			//购买的权限
			if roles_version != nil {
				_menu, _ := Menu_Com.Where("isbase", 0).OrWhereIn("id", roles_version_int).Where("type", "menu").Order("weigh asc").Get()
				menu = _menu
			} else {
				_menu, _ := Menu_Com.Where("isbase", 0).Order("weigh asc").Get()
				menu = _menu
			}
		} else {
			_menu, _ := Menu_Com.WhereIn("id", ids_arr_int).Order("weigh asc").Get()
			menu = _menu
		}
		menulist := GetTreeArray(menu, 0)
		//返回登录信息
		context.JSON(200, gin.H{
			"code":    0,
			"msg":     "查找用户数据成功！",
			"info":    userinfo,
			"role":    role,
			"menu":    menulist,
			"gruopid": group_ids,
			"time":    time.Now().Unix(),
		})
	} else {
		results.Failed(context, "对不起！您没有访问权限", err)
		return
	}
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
