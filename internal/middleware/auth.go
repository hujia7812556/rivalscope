// Package middleware 提供HTTP 中间件。
package middleware

import (
	"github.com/gin-gonic/gin"

	"rivalscope/internal/auth"
	"rivalscope/internal/dto"
)

// JWT 鉴权中间件:token 从 Authorization: Bearer <token> 读取,
// 缺失时回退读取 query 参数 accessToken(便于直接拼链接下载的场景)。
// 校验通过后将 userId 写入请求上下文。
func JWT(j *auth.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("accessToken")
		}
		if token == "" {
			c.AbortWithStatusJSON(401, dto.Response{
				Code: 401, Message: "未登录或登录已过期", Data: gin.H{},
			})
			return
		}
		claims, err := j.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, dto.Response{
				Code: 401, Message: "未登录或登录已过期", Data: gin.H{},
			})
			return
		}
		c.Set("userId", claims.UserId)
		c.Set("username", claims.Username)
		c.Set("nickname", claims.Nickname)
		c.Next()
	}
}
