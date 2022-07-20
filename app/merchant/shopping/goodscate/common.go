package goodscate

import (
	"basegin/app/model"

	"github.com/gohouse/gorose/v2"
)

func DB() gorose.IOrm {
	return model.DB.NewOrm()
}

//tool-获取树状数组
func GetTreeArray(num []gorose.Data, pid int64, itemprefix string) []gorose.Data {
	childs := ToolFar(num, pid) //获取pid下的所有数据
	var chridnum []gorose.Data
	if childs != nil {
		var number int = 1
		var total int = len(childs)
		for _, v := range childs {
			j := ""
			k := ""
			if number == total {
				j += "└"
				k = ""
				if itemprefix != "" {
					k = "&nbsp;"
				}

			} else {
				j += "├"
				k = ""
				if itemprefix != "" {
					k = "│"
				}
			}
			spacer := ""
			if itemprefix != "" {
				spacer = itemprefix + j
			}
			v["spacer"] = spacer
			v["childlist"] = GetTreeArray(num, v["id"].(int64), itemprefix+k+"&nbsp;")
			chridnum = append(chridnum, v)
			number++
		}
	}
	return chridnum
}

//2.将getTreeArray的结果返回为二维数组
func GetTreeList_txt(data []gorose.Data, field string) []gorose.Data {
	var midleArr []gorose.Data
	for _, v := range data {
		var childlist []gorose.Data
		if _, ok := v["childlist"]; ok {
			childlist = v["childlist"].([]gorose.Data)
		} else {
			childlist = make([]gorose.Data, 0)
		}
		delete(v, "childlist")
		v[field+"_txt"] = v["spacer"].(string) + " " + v[field+""].(string)
		if len(childlist) > 0 {
			v["haschild"] = 1
		} else {
			v["haschild"] = 0
		}
		if _, ok := v["id"]; ok {
			midleArr = append(midleArr, v)
		}
		if len(childlist) > 0 {
			newarr := GetTreeList_txt(childlist, field)
			midleArr = ArrayMerge(midleArr, newarr)
		}
	}
	return midleArr
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

//数组拼接
func ArrayMerge(ss ...[]gorose.Data) []gorose.Data {
	n := 0
	for _, v := range ss {
		n += len(v)
	}
	s := make([]gorose.Data, 0, n)
	for _, v := range ss {
		s = append(s, v...)
	}
	return s
}

//tool-获取树状数组
func GetTreeArray_only(num []gorose.Data, pid int64) []gorose.Data {
	childs := ToolFar(num, pid) //获取pid下的所有数据
	var chridnum []gorose.Data
	if childs != nil {
		for _, v := range childs {
			v["children"] = GetTreeArray_only(num, v["id"].(int64))
			chridnum = append(chridnum, v)
		}
	}
	return chridnum
}
