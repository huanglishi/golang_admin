package schedul

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
)

//添加模板
func Pushtpl_v1(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	f_id := parameter["id"].(float64)
	if f_id == 0 {
		parameter["accountID"] = user.Accountid
		addId, err := DB().Table("m_workermanage_staff_schedul").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("m_workermanage_staff_schedul").
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

//获取排班时模板
func Getschedultpl_v1(context *gin.Context) {
	dotype := context.DefaultQuery("dotype", "no")
	tid := context.DefaultQuery("tid", "")
	uid := context.DefaultQuery("uid", "0")
	//post数据
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if tid == "" {
		_tid, _ := DB().Table("m_workermanage_staff_user").Where("id", uid).Fields("schedultpl_id").First()
		tid = strconv.FormatInt(_tid["schedultpl_id"].(int64), 10)
	}
	if dotype == "yes" { //更新用户模板
		DB().Table("m_workermanage_staff_user").
			Data(map[string]interface{}{"schedultpl_id": tid}).
			Where("id", uid).
			Update()
	}
	tpllist, _ := DB().Table("m_workermanage_staff_schedul_item").Where("pid", tid).Order("weigh asc").Get()
	//处理时段排班状态
	for _, val := range tpllist {
		staff_user, _ := DB().Table("m_workermanage_staff_user").Where("id", uid).Fields("schedul_text").First()
		val["schedul_text"] = staff_user["schedul_text"]
		//处理排班数据
		tddata := []map[string]interface{}{}
		for _, cval := range parameter["timefield"].([]interface{}) {
			//查找数据
			schedul, _ := DB().Table("m_workermanage_staff_user_schedul").Where("uid", uid).Where("datetime", cval).Where("timeframe_id", val["id"]).Fields("id,ischeck,number").First()
			if schedul == nil {
				tddata = append(tddata, map[string]interface{}{"ischeck": false, "number": 0, "datetime": cval})
			} else {
				ischeck := false
				if schedul["ischeck"].(int64) == 1 {
					ischeck = true
				}
				tddata = append(tddata, map[string]interface{}{"ischeck": ischeck, "number": schedul["number"], "datetime": cval})
			}
		}
		val["tddata"] = tddata
	}
	results.Success(context, "获取数据成功", tpllist, nil)
}

//获取数据
func Gettpl_v1(context *gin.Context) {
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	defaultlist, _ := DB().Table("m_workermanage_staff_schedul_item").Where("pid", 0).Order("weigh asc").Get()
	//自定义数据
	customlist, _ := DB().Table("m_workermanage_staff_schedul").Where("accountID", user.Accountid).Fields("id,title").Get()
	for _, val := range customlist {
		list, _ := DB().Table("m_workermanage_staff_schedul_item").Where("pid", val["id"]).Order("weigh asc").Get()
		for _, cval := range list {
			cval["editable"] = false
		}
		val["timeframe"] = list
	}
	results.Success(context, "获取数据成功", map[string]interface{}{"default": defaultlist, "custom": customlist}, nil)
}

//保存时段
func Pushtimeframe_v1(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	f_id := parameter["id"].(float64)
	if f_id == 0 {
		addId, err := DB().Table("m_workermanage_staff_schedul_item").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			if addId != 0 {
				DB().Table("m_workermanage_staff_schedul_item").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("m_workermanage_staff_schedul_item").
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

//删除时段
func Deltimeframe_v1(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("m_workermanage_staff_schedul_item").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

//删除模板
func Deltpl_v1(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("m_workermanage_staff_schedul").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		//同时删除模板下的数据
		DB().Table("m_workermanage_staff_schedul_item").Where("pid", parameter["id"]).Delete()
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

//设置或取消排班
func Setschedul_v1(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter["ischeck"] == true { //设置
		schedul, _ := DB().Table("m_workermanage_staff_user_schedul").Where("uid", parameter["uid"]).Where("datetime", parameter["datetime"]).Where("timeframe_id", parameter["timeframe_id"]).Fields("id").First()
		var err error
		if schedul == nil {
			_, _err := DB().Table("m_workermanage_staff_user_schedul").Data(parameter).InsertGetId()
			err = _err
		} else {
			_, _err := DB().Table("m_workermanage_staff_user_schedul").
				Data(parameter).
				Where("id", schedul["id"]).
				Update()
			err = _err
		}
		if err != nil {
			results.Failed(context, "设置失败", err)
		} else {
			results.Success(context, "设置成功！", nil, nil)
		}
	} else { //取消
		schedul, _ := DB().Table("m_workermanage_staff_user_schedul").Where("uid", parameter["uid"]).Where("datetime", parameter["datetime"]).Where("timeframe_id", parameter["timeframe_id"]).Fields("id").First()
		if schedul == nil {
			results.Failed(context, "无需取消", nil)
		} else {
			res, err := DB().Table("m_workermanage_staff_user_schedul").
				Data(parameter).
				Where("id", schedul["id"]).
				Update()
			if err != nil {
				results.Failed(context, "取消失败", err)
			} else {
				results.Success(context, "取消成功！", res, nil)
			}
		}
	}
}
