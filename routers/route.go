package routers

import (
	rmold "basegin/routers/mold"
	"basegin/utils/Toolconf"
	"basegin/utils/handler"
	utils "basegin/utils/tool"
	"net/http"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	R *gin.Engine
)

func LimitHandler(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "您的操作太频繁，请稍后再试！",
				"data": nil,
			})
			c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}
func init() {

	R = gin.Default()
	gin.SetMode(gin.DebugMode) //ReleaseMode 为方便调试，Gin 框架在运行的时候默认是debug模式，在控制台默认会打印出很多调试日志，上线的时候我们需要关闭debug模式，改为release模式。
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	R.MaxMultipartMemory = 8 << 20 // 8 MiB

	//1.限流rate-limit 中间件
	lmt := tollbooth.NewLimiter(100, nil)
	lmt.SetMessage("")
	R.Use(LimitHandler(lmt))

	//2.部署vue项目
	// R.LoadHTMLGlob("viewst/*.html")              // 添加入口index.html
	R.Static("/resource", "./resource") // 附件
	R.StaticFile("/favicon.ico", "./resource/favicon.ico")
	R.Static("/views", "./views") // 管理后台
	// R.StaticFile("/viewst/", "viewst/index.html") //前端接口

	//3.跨域访问
	str_arr := strings.Split(Toolconf.AppConfig.String("allowurl"), `,`)
	R.Use(cors.New(cors.Config{
		AllowOrigins: str_arr,
		// AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "*"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Access-Token", "X-Requested-With", "Accept-Ranges", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Credentials", "Access-Control-Expose-Headers"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	//4.注意 Recover 要尽量放在第一个被加载
	R.Use(handler.Recover)

	//5.验证token
	R.Use(utils.JwtVerify)
	//6.找不到路由
	R.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		c.JSON(400, gin.H{"code": 400, "message": "您" + method + "请求地址：" + path + "不存在！"})
	})
	//7.路由配置文件
	rmold.Api(R)
	rmold.Apic(R)
	rmold.Apim(R)
	rmold.ApiWxs(R)
	rmold.Apiwxapp(R)
}

// func api() {
// 	auth := R.Group("/api")
// 	{
// 		// 模拟添加一条Policy策略
// 		auth.POST("acs", func(c *gin.Context) {
// 			APIResponse.C = c
// 			subject := "tom"
// 			object := "/api/routers"
// 			action := "POST"
// 			cacheName := subject + object + action
// 			result := ACS.Enforcer.AddPolicy(subject, object, action)
// 			if result {
// 				// 清除缓存
// 				_ = Cache.GlobalCache.Delete(cacheName)
// 				APIResponse.Success("add success")
// 			} else {
// 				APIResponse.Error("add fail")
// 			}
// 		})
// 		// 模拟删除一条Policy策略
// 		auth.DELETE("acs/:id", func(context *gin.Context) {
// 			APIResponse.C = context
// 			result := ACS.Enforcer.RemovePolicy("tom", "/api/routers", "POST")
// 			if result {
// 				// 清除缓存 代码省略
// 				APIResponse.Success("delete Policy success")
// 			} else {
// 				APIResponse.Error("delete Policy fail")
// 			}
// 		})
// 		// 获取路由列表
// 		auth.POST("/routers", middleware.Privilege(), func(c *gin.Context) {
// 			type data struct {
// 				Method string `json:"method"`
// 				Path   string `json:"path"`
// 			}
// 			var datas []data
// 			routers := R.Routes()
// 			for _, v := range routers {
// 				var temp data
// 				temp.Method = v.Method
// 				temp.Path = v.Path
// 				datas = append(datas, temp)
// 			}
// 			APIResponse.C = c
// 			APIResponse.Success(datas)
// 			return
// 		})
// 	}
// 	// 定义路由组
// 	user := R.Group("/api/v1")
// 	// 使用访问控制中间件
// 	user.Use(middleware.Privilege())
// 	{

// 		user.POST("user", func(c *gin.Context) {
// 			c.JSON(200, gin.H{"code": 200, "message": "user add success"})
// 		})
// 		user.DELETE("user/:id", func(c *gin.Context) {
// 			id := c.Param("id")
// 			c.JSON(200, gin.H{"code": 200, "message": "user delete success " + id})
// 		})
// 		user.PUT("user/:id", func(c *gin.Context) {
// 			id := c.Param("id")
// 			c.JSON(200, gin.H{"code": 200, "message": "user update success " + id})
// 		})
// 		user.GET("user/:id", func(c *gin.Context) {
// 			id := c.Param("id")
// 			c.JSON(200, gin.H{"code": 200, "message": "user Get success " + id})
// 		})
// 		user.GET("test", func(c *gin.Context) {
// 			c.JSON(200, gin.H{"code": 200, "message": "测试函数2"})
// 		})
// 	}
// }
