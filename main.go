package main

import (
	"basegin/app/model"
	. "basegin/routers"
	"basegin/utils/Toolconf"
	"fmt"
	"runtime"
	"strconv"

	_ "github.com/go-sql-driver/mysql" //只执行github.com/go-sql-driver/mysql的init函数
	"github.com/gohouse/gorose/v2"     //数据库操作
)

var db *gorose.Engin

func InitDB() {
	fmt.Println("数据库初始化中...")
	var err error
	db, err = gorose.Open(&gorose.Config{Driver: "mysql", Dsn: Toolconf.AppConfig.String("db.user") + ":" + Toolconf.AppConfig.String("db.password") + "@tcp(" + Toolconf.AppConfig.String("db.host") + ":" + Toolconf.AppConfig.String("db.port") + ")/" + Toolconf.AppConfig.String("db.name") + "?charset=utf8mb4&parseTime=true&loc=Local", SetMaxOpenConns: 100, SetMaxIdleConns: 10})
	if err != nil {
		fmt.Println("链接数据库错误，请检查数据库链接！")
	}
	model.DB = db
}
func main() {
	//多核并行任务
	cpu_num, _ := strconv.Atoi(Toolconf.AppConfig.String("cpunum"))
	mycpu := runtime.NumCPU()
	if cpu_num > mycpu { //如果配置cpu核数大于当前计算机核数，则等当前计算机核数
		cpu_num = mycpu
	}
	if cpu_num > 0 {
		fmt.Printf("当前计算机核数: %v个,调用：%v个\n", mycpu, cpu_num)
		runtime.GOMAXPROCS(cpu_num)
	} else {
		fmt.Printf("当前计算机核数: %v个,调用：%v个\n", mycpu, mycpu)
		runtime.GOMAXPROCS(mycpu)
	}
	InitDB()
	// Run("里面不指定端口号默认为8088")
	R.Run(Toolconf.AppConfig.String("httpport"))
}
