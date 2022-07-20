package wxmenu

import (
	"basegin/app/model"
	"encoding/json"
	"strconv"

	"github.com/gohouse/gorose/v2"
)

func DB() gorose.IOrm {
	return model.DB.NewOrm()
}

/**
* 更具当前用户id返回主账号id
 */
func GetMainUID(uid int64) int64 {
	var bid int64
	pid, _ := DB().Table("merchant_user").Where("id", uid).Value("pid")
	var zon int64 = 0
	pid_ := pid.(int64)
	if pid_ != zon {
		f_pid, _ := DB().Table("merchant_user").Where("id", pid_).Value("pid")
		f_pid_ := f_pid.(int64)
		if f_pid_ != zon { //如果还不是主账号递归
			GetMainUID(f_pid_)
		} else {
			bid = pid_
		}
	} else {
		bid = uid
	}
	return bid
}

// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}
