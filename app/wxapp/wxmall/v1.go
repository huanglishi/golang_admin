package wxmall

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
	//获取店铺信息
	setting, _err := DB().Table("m_shopping_setting").Where("accountID", user.Accountid).Fields("id,name,tel,address,location,location_name,logo,door_img,notice_msg,des").First()
	//获取分类
	catedata, _err := DB().Table("m_shopping_goods_cate").Where("accountID", user.Accountid).Where("status", 0).Where("show_nav", 1).Fields("id,name").Order("weigh asc,id desc").Get()
	results.Success(context, "获取基础数据", map[string]interface{}{"setting": setting, "catedata": catedata}, _err)
}

//获取产品列表
func Getlist_v1(context *gin.Context) {
	cid := context.DefaultQuery("cid", "0")
	list, _err := DB().Table("m_shopping_goods_content").Where("status", 0).Where("cid", cid).Fields("id,title,price,thumb,des,sales").Order("weigh asc,id desc").Get()
	if _err != nil {
		results.Failed(context, "获取产品列表失败！", _err)
	} else {
		results.Success(context, "获取产品列表", list, nil)
	}
}

//获取商品详情页数据
func Getdetail_v1(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	//1.获取商品详情
	goodsdata, _err := DB().Table("m_shopping_goods_content").Where("id", id).Fields("id,title,price,original_price,stock_num,unit,des,thumb,photo,content,skulist,properties,propertieslist").First()
	if _err != nil {
		results.Failed(context, "获取商品详情失败！", _err)
	} else {
		skudata, _ := DB().Table("m_shopping_goods_content_sku").Where("product_id", id).Get()
		for _, val := range skudata {
			var parameter interface{}
			json.Unmarshal([]byte(val["rule_table_key"].(string)), &parameter)
			val["rule_table_key"] = parameter
		}
		goodsdata["sku"] = skudata
		//获取购物车产品数
		getuser, _ := context.Get("user") //当前用户
		user := getuser.(*utils.UserClaims)
		cartnum, _ := DB().Table("m_shopping_cart").Where("uid", user.ID).Count()
		results.Success(context, "获取商品详情", map[string]interface{}{"goodsdata": goodsdata, "cartnum": cartnum}, nil)
	}
}

//加入购车
func Addcart_v1(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	cartdata, _err := DB().Table("m_shopping_cart").Where("uid", user.ID).Where("pid", parameter["pid"]).Where("sku_id", parameter["sku_id"]).Fields("id,number,price").First()
	if _err != nil {
		results.Failed(context, "加入购车失败", _err)
	} else {
		if cartdata != nil {
			res2, err := DB().Table("m_shopping_cart").Where("id", cartdata["id"]).Data(map[string]interface{}{"number": cartdata["number"].(int64) + 1, "price": cartdata["price"].(int64) + parameter["price"].(int64)}).Update()
			if err != nil {
				results.Failed(context, "加入购车失败", err)
			} else {
				results.Success(context, "加入购车成功！", 0, res2)
			}
		} else {
			parameter["accountID"] = user.Accountid
			parameter["uid"] = user.ID
			addId, err := DB().Table("m_shopping_cart").Data(parameter).InsertGetId()
			if err != nil {
				results.Failed(context, "加入购车失败", err)
			} else {
				results.Success(context, "加入购车成功！", addId, nil)
			}
		}
	}
}

//获取购物车数据
func Getcart_v1(context *gin.Context) {
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	list, _err := DB().Table("m_shopping_cart").Where("uid", user.ID).Fields("id,pid,price,sku_id,properties,number,messages").Order("id desc").Get()
	if _err != nil {
		results.Failed(context, "获取购物车数据失败！", _err)
	} else {
		for _, val := range list {
			goodsdata, _ := DB().Table("m_shopping_goods_content").Where("id", val["pid"]).Fields("id,title,price,stock_num,thumb").First()
			val["goodsdata"] = goodsdata
			skudata, _ := DB().Table("m_shopping_goods_content_sku").Where("id", val["sku_id"]).Fields("id,image,price,stock_num,val").First()
			val["skudata"] = skudata
		}
		results.Success(context, "获取购物车数据", list, nil)
	}
}

//删除购物车
func Delcart_v1(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("m_shopping_cart").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}
