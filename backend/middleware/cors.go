package middleware

import (
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有来源访问（开发环境）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		
		// 允许携带认证信息（如 cookies）
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		
		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", 
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, "+
			"Authorization, Accept, Origin, Cache-Control, X-Requested-With, "+
			"Access-Control-Request-Headers, Access-Control-Request-Method")
		
		// 允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", 
			"POST, OPTIONS, GET, PUT, DELETE, PATCH")
		
		// 预检请求缓存时间
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		
		// 如果是预检请求，直接返回 204
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
} 