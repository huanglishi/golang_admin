package middleware

import (
	"basegin/utils/APIResponse"
	"basegin/utils/Cache"

	"github.com/gin-gonic/gin"
)

func Privilege() gin.HandlerFunc {
	return func(c *gin.Context) {
		APIResponse.C = c
		var userName = c.GetHeader("userName")
		if userName == "" {
			APIResponse.Error("参数缺失")
			c.Abort()
			return
		}
		path := c.Request.URL.Path
		method := c.Request.Method
		cacheName := userName + path + method
		// 从缓存中读取&判断
		entry, err := Cache.GlobalCache.Get(cacheName)
		if err == nil && entry != nil {
			if string(entry) == "true" {
				c.Next()
			} else {
				APIResponse.Error("从缓存中读取：无访问权限")
				c.Abort()
				return
			}
		} else {
		}
	}
}
