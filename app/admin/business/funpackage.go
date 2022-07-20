package business

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

//获取套餐市场
func Getmarketdata(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	var cid []interface{}
	industry, _ := DB().Table("merchant_user").Where("id", id).Value("industry")
	if industry != nil {
		cid = append(cid, industry)
	}
	module_cate, _ := DB().Table("admin_dev_module_cate").Where("status", 0).Where("type", 1).Pluck("id")
	for _, v := range module_cate.([]interface{}) {
		cid = append(cid, v)
	}
	// _pageNo := context.DefaultQuery("pageNo", "1")
	// _pageSize := context.DefaultQuery("pageSize", "10")
	list, _ := DB().Table("admin_dev_module_content").WhereIn("cid", cid).Order("id asc").Get()
	for _, val := range list {
		version, _ := DB().Table("admin_dev_module_version").Where("mcid", val["id"]).Fields("id,api,name,price").Order("weigh asc").Get()
		val["version"] = version
		if len(version) > 0 {
			val["price"] = version[0]["price"]
			val["vid"] = version[0]["id"]
		}
		catename, _ := DB().Table("admin_dev_module_cate").Where("id", val["cid"]).Value("name")
		val["catename"] = catename

	}
	results.Success(context, "套餐市场列表", list, nil)
}

//获取版本信息
func Versionview(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	list, _ := DB().Table("admin_dev_module_version").Where("mcid", id).Fields("id,api,name,des").Order("weigh asc").Get()
	results.Success(context, "获取版本信息列表", list, nil)
}

//保存购物车数据
func Savecar(context *gin.Context) {
	//获取post传过来的data
	aid := context.DefaultQuery("aid", "0")
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter []map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	//批量提交
	save_arr := []map[string]interface{}{}
	for _, val := range parameter {
		choosepack, _ := DB().Table("admin_business_choosepack").Where("accountID", aid).Where("module_id", val["id"]).Fields("id").Distinct().First()
		if choosepack == nil {
			marr := map[string]interface{}{"uid": user.ID, "accountID": aid, "module_id": val["id"], "vid": val["vid"], "price": val["price"], "createtime": time.Now().Unix()}
			save_arr = append(save_arr, marr)
		}
	}
	res, err := DB().Table("admin_business_choosepack").Data(save_arr).Insert()
	if err != nil && len(save_arr) > 0 {
		results.Failed(context, "保存购物车数据失败", err)
	} else {
		results.Success(context, "保存购物车数据成功", res, nil)
	}

}

//获取已选套餐
func Choosepackage(context *gin.Context) {
	aid := context.DefaultQuery("aid", "0")
	list, _ := DB().Table("admin_dev_module_content a").LeftJoin("admin_business_choosepack b on a.id = b.module_id").Where("b.accountID", aid).Fields("a.id,a.cid,a.title,a.des,a.image,a.weigh,b.accountID,b.vid,b.price,b.createtime,b.ispay,b.validtime,b.buytime").AddFields("b.id as bid").Get()
	for _, val := range list {
		version, _ := DB().Table("admin_dev_module_version").Where("id", val["vid"]).Fields("id,api,name,price").Distinct().First()
		val["version"] = version
		catename, _ := DB().Table("admin_dev_module_cate").Where("id", val["cid"]).Value("name")
		val["catename"] = catename
	}
	//统计
	paynum, _ := DB().Table("admin_business_choosepack").Where("accountID", aid).Where("ispay", 1).Count()
	paymoney, _ := DB().Table("admin_business_choosepack").Where("accountID", aid).Where("ispay", 1).Sum("price")
	nopaynum, _ := DB().Table("admin_business_choosepack").Where("accountID", aid).Where("ispay", 0).Count()
	nopaymoney, _ := DB().Table("admin_business_choosepack").Where("accountID", aid).Where("ispay", 0).Sum("price")
	oldTime := time.Now().AddDate(0, 0, 5)
	duesoon, _ := DB().Table("admin_business_choosepack").Where("accountID", aid).Where("validtime", "!=", 0).Where("validtime", "<", oldTime.Unix()).Where("validtime", ">", time.Now().Unix()).Count()
	expired, _ := DB().Table("admin_business_choosepack").Where("accountID", aid).Where("validtime", "!=", 0).Where("validtime", "<", time.Now().Unix()).Count()
	var statistic = map[string]interface{}{"paynum": paynum, "paymoney": paymoney, "nopaynum": nopaynum, "nopaymoney": nopaymoney, "duesoon": duesoon, "expired": expired}
	results.Success(context, "获取已选套餐", list, statistic)
}

//删除选择的数据-批量
func Delchoose(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res, err := DB().Table("admin_business_choosepack").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else if res == 0 {
		results.Failed(context, "没有找到删除数据", err)
	} else {
		results.Success(context, "删除成功！", res, nil)
	}
	context.Abort()
	return
}

// 更新套餐数据
func Upchoose(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	bid, _ := parameter["bid"]
	delete(parameter, "bid")
	if parameter["validtime"] != nil {
		validtime, _ := time.ParseInLocation("2006/01/02 ", parameter["validtime"].(string), time.Local)
		parameter["validtime"] = validtime.Unix()
	} else {
		parameter["validtime"] = 0
	}
	res2, err := DB().Table("admin_business_choosepack").Where("id", bid).Data(parameter).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		msg := "更新成功！"
		if res2 == 0 {
			msg = "暂无数据更新"
		} else { //修改权限

		}
		results.Success(context, msg, res2, nil)
	}

}
