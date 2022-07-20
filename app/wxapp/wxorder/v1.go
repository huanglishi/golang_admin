package wxorder

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

//加入订单
func Addorder_v1(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter []map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)

	oid, err := DB().Table("m_shopping_order").Data(map[string]interface{}{"accountID": user.Accountid, "uid": user.ID, "createtime": time.Now().Unix()}).InsertGetId()
	if err != nil {
		results.Failed(context, "创建订单失败", err)
	} else {
		//添加购物清单
		save_arr := []map[string]interface{}{}
		for _, narr := range parameter {
			narr["oid"] = oid
			properties, _ := utils.JsonMarshalNoSetEscapeHTML(narr["properties"])
			narr["properties"] = properties
			save_arr = append(save_arr, narr)
		}
		DB().Table("m_shopping_orderlist").Where("oid", oid).Delete()
		if len(save_arr) > 0 {
			DB().Table("m_shopping_orderlist").Data(save_arr).Insert()
		}
		results.Success(context, "创建订单成功！", oid, nil)
	}

}
