package docappointment

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

//获取获取基础数据
func Getbase_v1(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	catedata, _err := DB().Table("m_workermanage_staff_group").Where("accountID", user.Accountid).Where("status", 0).Where("show_nav", 1).Fields("id,name").Order("weigh asc,id desc").Get()

	results.Success(context, "获取基础数据", map[string]interface{}{"catedata": catedata}, _err)
}

//获取医生列表
func Getlist_v1(context *gin.Context) {
	cid := context.DefaultQuery("cid", "0")
	_cate_ids, _ := DB().Table("m_workermanage_staff_group").Where("pid", cid).Pluck("id")
	cate_ids := _cate_ids.([]interface{})
	cate_ids = append(cate_ids, cid)
	workerids, _ := DB().Table("m_workermanage_staff_group_bind").WhereIn("group_id", cate_ids).Pluck("worker_id")
	_workerids := workerids.([]interface{})
	if len(_workerids) > 0 {
		list, _err := DB().Table("m_workermanage_staff_user").Where("status", 0).WhereIn("id", _workerids).Fields("id,name,headimgurl,position,professional,skilled,remark").Order("weigh asc,id desc").Get()
		if _err != nil {
			results.Failed(context, "获取医生列表失败！", _err)
		} else {
			results.Success(context, "获取医生列表", list, nil)
		}
	} else {
		results.Success(context, "获取医生列表", _workerids, nil)
	}
}

//获取医生详情页
func Getdetail_v1(context *gin.Context) {
	uid := context.DefaultQuery("id", "0")
	//post数据
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//1.获取医生信息
	docdata, _err := DB().Table("m_workermanage_staff_user").Where("id", uid).Fields("id,name,headimgurl,position,professional,remark,skilled,accepts,consult,evaluate").First()
	if _err != nil {
		results.Failed(context, "获取医生详情页失败！", _err)
	} else {
		//2.获取排班数据
		//2.1获取排班模板
		tid, _ := DB().Table("m_workermanage_staff_user").Where("id", uid).Value("schedultpl_id")
		tpllist, _ := DB().Table("m_workermanage_staff_schedul_item").Where("pid", tid).Order("weigh asc").Get()
		//排班状态
		dataitem := parameter["weekday"].([]interface{})
		for _, weekday := range dataitem {
			weekday_arr := weekday.(map[string]interface{})
			tddata := []map[string]interface{}{}
			for _, val := range tpllist {
				//查找数据
				schedul, _ := DB().Table("m_workermanage_staff_user_schedul").Where("uid", uid).Where("datetime", weekday_arr["datetime"]).Where("timeframe_id", val["id"]).Fields("id,ischeck,number,use_number").First()
				if schedul == nil {
					tddata = append(tddata, map[string]interface{}{"ischeck": false, "number": 0, "datetime": weekday_arr["datetime"]})
				} else {
					ischeck := false
					if schedul["ischeck"].(int64) == 1 {
						ischeck = true
					}
					tddata = append(tddata, map[string]interface{}{"ischeck": ischeck, "number": schedul["number"], "use_number": schedul["use_number"], "datetime": weekday_arr["datetime"], "start_time": val["start_time"], "end_time": val["end_time"]})
				}
			}
			weekday_arr["tddata"] = tddata
		}
		//查找关注状态
		docdata["subscribe"] = 0
		getuser, _ := context.Get("user") //当前用户
		user := getuser.(*utils.UserClaims)
		recordtdata, _ := DB().Table("m_dotag_record").Where("uid", user.ID).Where("type", 1).Where("obj_id", uid).Fields("id").First()
		if recordtdata != nil {
			docdata["subscribe"] = recordtdata["id"]
		}
		results.Success(context, "获取基础数据", map[string]interface{}{"docdata": docdata, "tpllist": tpllist, "trdata": dataitem}, parameter["timefield"])
	}
}

//获取关注的医生
func Getmydoctor_v1(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	doc_ids, _err := DB().Table("m_dotag_record").Where("uid", user.ID).Where("type", 1).Pluck("obj_id")
	list, _err := DB().Table("m_workermanage_staff_user").Where("status", 0).WhereIn("id", doc_ids.([]interface{})).Fields("id,name,headimgurl,position,professional,skilled,remark").Order("weigh asc,id desc").Get()
	if _err != nil {
		results.Failed(context, "获取医生列表失败！", _err)
	} else {
		results.Success(context, "获取医生列表", list, nil)
	}
}
