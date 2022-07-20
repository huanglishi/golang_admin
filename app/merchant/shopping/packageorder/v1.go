package packageorder

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
)

//获取内容数据
func Cetlist(context *gin.Context) {
	title := context.DefaultQuery("title", "")
	keywords := context.DefaultQuery("keywords", "")
	createtime := context.DefaultQuery("createtime", "")
	_pageNo := context.DefaultQuery("pageNo", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(_pageNo)
	pageSize, _ := strconv.Atoi(_pageSize)
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)

	Fdb := DB().Table("m_shopping_packageorder").Where("accountID", user.Accountid)
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
	getfield := "*"
	list, err := Fdb.Fields(getfield).Limit(pageSize).Page(pageNo).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "查找失败", err)
	} else {
		// 统计数据
		var totalCount int64
		totalCount, _ = DB().Table("m_shopping_packageorder").Where("accountID", user.Accountid).Count()
		_pageSize := int64(pageSize)
		totalPage := totalCount / _pageSize
		// for _, val := range list {
		// 	orderlist, _ := DB().Table("m_shopping_orderlist").Where("pid", val["id"]).Order("id desc").Get()
		// 	val["orderlist"] = orderlist
		// }
		results.Success(context, "查找成功！", map[string]interface{}{
			"pageNo":     pageNo,
			"pageSize":   pageSize,
			"totalCount": totalCount,
			"totalPage":  totalPage,
			"data":       list}, nil)
	}
}

//删除内容-批量
func Delcontent_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("m_shopping_packageorder").WhereIn("id", ids.([]interface{})).Delete()
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
	res2, err := DB().Table("m_shopping_packageorder").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
