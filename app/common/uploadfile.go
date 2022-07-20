package common

import (
	"basegin/utils/results"
	utils "basegin/utils/tool"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//上传单文件
func OneFile(context *gin.Context) {
	// 单个文件
	file, err := context.FormFile("file")
	if err != nil {
		results.Failed(context, "获取数据失败，", err)
		return
	}
	sha1_str := md5Str(file.Filename)
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	attachment, _ := DB().Table("attachment").Where("uid", user.ID).Where("sha1", sha1_str).Fields("title,url").First()
	if attachment != nil { //文件是否已经存在
		context.JSON(200, gin.H{
			"uid":      sha1_str,
			"name":     attachment["title"],
			"status":   "done",
			"url":      attachment["url"],
			"response": "文件已上传",
			"time":     time.Now().Unix(),
		})
		context.Abort()
		return
	}
	file_path := fmt.Sprintf("%s%s%s", "resource/uploads/", time.Now().Format("20060102"), "/")
	//如果没有filepath文件目录就创建一个
	if _, err := os.Stat(file_path); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(file_path, os.ModePerm)
		}
	}
	//上传到的路径
	//path := 'resource/uploads/20060102150405test.xlsx'
	file_Filename := fmt.Sprintf("%v%s%s", time.Now().Unix(), "_", file.Filename) // 文件名格式 自己可以改 建议保证唯一性
	path := file_path + file_Filename                                             //路径+文件名上传
	// 上传文件到指定的目录
	err = context.SaveUploadedFile(file, path)
	if err != nil {
		context.JSON(200, gin.H{
			"uid":      sha1_str,
			"name":     file.Filename,
			"status":   "error",
			"response": "上传失败",
			"time":     time.Now().Unix(),
		})
	} else {
		//保存数据
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		Insertdata := map[string]interface{}{
			"accountID":  user.Accountid,
			"uid":        user.ID,
			"sha1":       sha1_str,
			"title":      file.Filename,
			"url":        path,
			"storage":    dir + strings.Replace(path, "/", "\\", -1),
			"uploadtime": time.Now().Unix(),
			"updatetime": time.Now().Unix(),
			"filesize":   file.Size,
			"mimetype":   file.Header["Content-Type"][0],
		}
		DB().Table("attachment").Data(Insertdata).Insert()
		context.JSON(200, gin.H{
			"uid":      sha1_str,
			"name":     file.Filename,
			"status":   "done",
			"url":      path,
			"thumb":    path,
			"response": "上传成功",
			// "file":     file.Header,
			"time": time.Now().Unix(),
		})
	}
}

//显示图片
func GetImage(context *gin.Context) {
	imageName := context.Query("url")
	context.File(imageName)
}

//显示图片base64
func Getimagebase(context *gin.Context) {
	imageName := context.Query("url")
	file, _ := ioutil.ReadFile(imageName)
	context.Writer.WriteString(string(file))
}
func md5Str(origin string) string {
	m := md5.New()
	m.Write([]byte(origin))
	return hex.EncodeToString(m.Sum(nil))
}

func Testpath(context *gin.Context) {
	log.Printf("测试调度: %v\n", time.Now().Unix())
	context.JSON(200, gin.H{
		"Code":     200,
		"response": "测试调度",
	})
}
