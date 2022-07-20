package edithtmltpl

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

//创建和获取文件夹及文件
func Getsitefile(context *gin.Context) {
	site_id := context.DefaultQuery("site_id", "0")
	file_path := fmt.Sprintf("%s%s%v%s", "website/public/template/site", "_", site_id, "/")
	//如果没有filepath文件目录就创建一个
	if _, err := os.Stat(file_path); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(file_path, os.ModePerm)
		}
	}
	flist, _ := GetAllFiles(file_path)
	nowsite := fmt.Sprintf("%s%s%v", "site", "_", site_id)
	results.Success(context, "创建和获取文件夹及文件", flist, nowsite)
}

func GetAllFiles(dirPth string) (files []string, err error) {
	fis, err := ioutil.ReadDir(filepath.Clean(filepath.ToSlash(dirPth)))
	if err != nil {
		return nil, err
	}
	for _, f := range fis {
		_path := filepath.Join(dirPth, f.Name())
		if f.IsDir() {
			fs, _ := GetAllFiles(_path)
			files = append(files, fs...)
		}
		files = append(files, _path)
	}
	return files, nil
}

//新建文件夹
func Newdir(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	//提交文件名
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	file_path := fmt.Sprintf("%s%s%v%s%s", "website/public/template/site", "_", parameter["site_id"], parameter["path"], parameter["name"])
	//如果没有filepath文件目录就创建一个
	if _, err := os.Stat(file_path); err != nil {
		if !os.IsExist(err) {
			err = os.MkdirAll(file_path, os.ModePerm)
			if err != nil {
				results.Failed(context, "新建文件夹失败", err)
				context.Abort()
				return
			}
		}
	}
	results.Success(context, "新建文件夹成功", user.Accountid, nil)
}

//新建文件
func Newfile(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	//提交文件名
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	dir_path := fmt.Sprintf("%s%s%v%s", "website/public/template/site", "_", parameter["site_id"], parameter["path"])
	file_path := fmt.Sprintf("%s%s%v%s%s", "website/public/template/site", "_", parameter["site_id"], parameter["path"], parameter["name"])
	//如果没有filepath文件目录就创建一个
	if _, err := os.Stat(dir_path); err != nil {
		if !os.IsExist(err) {
			err = os.MkdirAll(dir_path, os.ModePerm)
			if err != nil {
				results.Failed(context, "新建文件夹失败", err)
				context.Abort()
				return
			}
		}
	}
	//判断文件是否存在
	_, err := os.Lstat(file_path)
	if os.IsNotExist(err) {
		fp, err := os.Create(file_path) // 如果文件已存在，会将文件清空。
		// defer延迟调用
		defer fp.Close() //关闭文件，释放资源
		if err != nil {
			results.Failed(context, "新建文件失败", err)
			context.Abort()
			return
		}
	}
	results.Success(context, "新建文件夹成功", user.Accountid, nil)
}

//删除文件或文件夹
func Delfiles(context *gin.Context) {
	//提交文件名
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	file_path := fmt.Sprintf("%s%s%v%s", "website/public/template/site", "_", parameter["site_id"], parameter["path"])
	if parameter["dokey"] == "dir" { //删除文件夹
		err := os.RemoveAll(file_path)
		if err != nil {
			results.Failed(context, "删除文件夹失败", err)
			context.Abort()
			return
		}
	} else { //删除文件
		err := os.Remove(file_path)
		if err != nil {
			results.Failed(context, "删除文件失败", err)
			context.Abort()
			return
		}
	}
	results.Success(context, "删除成功", file_path, nil)
}

//读取文件中
func Readfile(context *gin.Context) {
	site_id := context.DefaultQuery("site_id", "0")
	path_str := context.DefaultQuery("path", "")
	filetype := context.DefaultQuery("filetype", "imgss")
	file_path := fmt.Sprintf("%s%s%v%s%s", "website/public/template/site", "_", site_id, "/", path_str)
	//如果没有filepath文件目录就创建一个
	if _, err := os.Stat(file_path); err != nil {
		if !os.IsExist(err) {
			results.Failed(context, "读取文件不存在！", err)
			context.Abort()
		}
	}
	if filetype == "img" {
		results.Success(context, "读取文件内容", file_path, filetype)
		context.Abort()
		return
	}
	bytes, err := ioutil.ReadFile(file_path)
	if err != nil {
		results.Failed(context, "读取文件内容失败", err)
		context.Abort()
	} else {
		results.Success(context, "读取文件内容", string(bytes), filetype)
	}
}

//写入文件
func Writefile(context *gin.Context) {
	//提交文件名
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	file_path := fmt.Sprintf("%s%s%v%s%s", "website/public/template/site", "_", parameter["site_id"], "/", parameter["path"])
	if _, err := os.Stat(file_path); err != nil {
		if !os.IsExist(err) {
			results.Failed(context, "写入文件不存在！", err)
			context.Abort()
			return
		}
	}
	f, err := os.OpenFile(file_path, os.O_WRONLY|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		results.Failed(context, "写入文件内容打开文件时失败", err.Error())
		context.Abort()
		return
	} else {
		_, err = f.Write([]byte(parameter["content"].(string)))
		if err != nil {
			results.Failed(context, "写入文件内容失败", err)
			context.Abort()
			return
		}
	}
	results.Success(context, "写入文件内容成功", file_path, nil)
}
