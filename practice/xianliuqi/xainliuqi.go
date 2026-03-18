package xianliuqi

import (
	"context"
	"fmt"
	"practice/redis"

	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ctx = context.Background()

const (
	ArticleLimitCount = 3
	ArticleLimitTime  = 1 * time.Hour
	CommentLimitTime  = 1 * time.Hour
	CommentCount      = 1
)

func GetRedisKey(userID uint, act string) string {
	return fmt.Sprintf("rate_limit:%s:%d", userID, act)
}
func CheckRateLimit(userID uint, act string) (bool, error) {
	var LimitCount int
	var LimitTime time.Duration
	switch act {
	case "article":
		LimitCount = ArticleLimitCount
		LimitTime = ArticleLimitTime
	case "comment":
		LimitCount = CommentCount
		LimitTime = CommentLimitTime
	default:
		return false, nil
	}
	key := GetRedisKey(userID, "article")
	count, err := redis.RedisClient.Incr(ctx, key).Result()
	if err != nil {
		zap.Error(err)
	}
	if count == 1 {
		redis.RedisClient.Expire(ctx, key, LimitTime)
	}
	if count > int64(LimitCount) {
		return true, nil
	}
	return false, nil
}
func ArticleLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exist := c.Get("userID")
		if !exist {
			c.JSON(200, gin.H{
				"code": 200,
			})
			c.Abort()
			return
		}
		uid, ok := userID.(uint)
		if !ok {
			c.JSON(200, gin.H{
				"code": 200,
			})
			c.Abort()
			return
		}
		exceed, err := CheckRateLimit(uint(uid), "article")
		if err != nil {
			c.JSON(200, gin.H{
				"code": 200,
			})
			c.Abort()
			return
		}
		if exceed {
			c.JSON(200, gin.H{
				"code": 200,
			})
			c.Abort()
			return
		}
	}
}
func CommentLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exist := c.Get("userID")
		if !exist {
			c.JSON(200, gin.H{
				"code": 200,
			})
			c.Abort()
			return
		}
		uid, ok := userID.(uint)
		if !ok {
			c.JSON(200, gin.H{
				"code": 200,
			})
			c.Abort()
			return
		}
		exceed, err := CheckRateLimit(uint(uid), "comment")
		if err != nil {
			c.JSON(200, gin.H{
				"code": 200,
			})
			c.Abort()
			return
		}
		if exceed {
			c.JSON(200, gin.H{
				"code": 200,
			})
			c.Abort()
			return
		}

	}
}
