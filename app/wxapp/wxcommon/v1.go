package wxcommon

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

//获取文章类详情
func Getarticledetail_v1(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	tablename := context.DefaultQuery("tablename", "")
	getfield := context.DefaultQuery("getfield", "")
	upfield := context.DefaultQuery("upfield", "")
	contentdata, _err := DB().Table(tablename).Where("id", id).Fields(getfield).First()
	if contentdata != nil {
		if upfield != "" {
			DB().Table(tablename).Where("id", id).Data(map[string]interface{}{upfield: contentdata[upfield].(int64) + 1}).Update()
		}
		results.Success(context, "获取详情", contentdata, _err)
	} else {
		results.Failed(context, "获取详情失败！", _err)
	}
}

//标记操作-关注，收藏
func Dotag_v1(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	recordtdata, _err := DB().Table("m_dotag_record").Where("uid", user.ID).Where("type", parameter["type"]).Where("obj_id", parameter["obj_id"]).Fields("id").First()
	if _err != nil {
		results.Failed(context, "操作失败", _err)
	} else {
		if recordtdata != nil {
			res2, err := DB().Table("m_dotag_record").Where("id", recordtdata["id"]).Delete()
			if err != nil {
				results.Failed(context, "取消失败", err)
			} else {
				results.Success(context, "取消成功！", 0, res2)
			}
		} else {
			parameter["uid"] = user.ID
			addId, err := DB().Table("m_dotag_record").Data(parameter).InsertGetId()
			if err != nil {
				results.Failed(context, "添加失败", err)
			} else {
				results.Success(context, "添加成功！", addId, nil)
			}
		}
	}
}

//获取附件
func Getfile_v1(context *gin.Context) {
	fileName := context.Query("faname")
	filePath := "resource/staticfile/" + fileName
	results.Failed(context, "添加失败", filePath)
	// context.File(filePath)
}
