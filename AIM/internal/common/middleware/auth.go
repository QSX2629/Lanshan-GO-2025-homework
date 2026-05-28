package middleware

import (
	"AIM/internal/common/error"
	"AIM/internal/common/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			error.Fail(c, 401, "请先登录")
			c.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := utils.ParseToken(token)
		if err != nil {
			error.Fail(c, 401, "登录已过期或token无效")
			c.Abort()
			return
		}
		// 存入上下文，业务层直接拿
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
