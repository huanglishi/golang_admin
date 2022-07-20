package cmdshell

import (
	"basegin/utils/results"
	"bytes"
	"fmt"
	"os/exec"

	"github.com/gin-gonic/gin"
)

//nginx 重新加载配置文件
func Reloadnginx(context *gin.Context) {
	// command := exec.Command("nssm restart nginx")
	// if runtime.GOOS == "windows" {
	// 	command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	// }
	// err := command.Run()

	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	//拼写命令
	// ml_str := "ping www.ynjiyuan.com"
	ml_str := "nssm restart nginx"
	cmd := exec.Command("cmd", "/C", ml_str)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out
	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	if err != nil {
		results.Failed(context, "执行失败", err)
		return
	}
	results.Success(context, "重新加载配置文件成功！", out.String(), err)

}

//错误处理函数
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
