package jwt

import (
	"context"
	"net/http"
	"practice/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("demo-secret-key")
var ctx = context.Background()

func GenerateToken(userID uint, username string, role string) (string, error) {
	claims := models.Claims{
		Username: username,
		UserID:   userID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
func ParseToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// AuthRequired 角色验证
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(200, gin.H{
				"error": "请登录",
			})
			c.Abort()
			return
		}
		claims, err := ParseToken(token)
		if err != nil || claims == nil {
			c.JSON(200, gin.H{
				"error": "token无效",
			})
			c.Abort()
			return
		}
		if claims.Role != "admin" {
			// 403 Forbidden 表示已登录但无权限，区别于401未登录
			c.JSON(http.StatusForbidden, gin.H{"error": "无管理员权限，拒绝访问"})
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}

}
