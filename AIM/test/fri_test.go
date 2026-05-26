package test

import (
	"AIM/internal/api"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAddFriend(t *testing.T) {
	r := gin.Default()
	// 直接注册包级函数
	r.POST("/friend/add", api.AddFriend)

	body := strings.NewReader(`{"friend_id":"test001"}`)
	req, _ := http.NewRequest("POST", "/friend/add", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	t.Log("✅ 添加好友测试完成")
}

func TestCreateFriendGroup(t *testing.T) {
	r := gin.Default()
	r.POST("/friend/group/create", api.CreateFriendGroup)

	body := strings.NewReader(`{"group_name":"我的好友"}`)
	req, _ := http.NewRequest("POST", "/friend/group/create", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	t.Log("✅ 创建好友分组测试完成")
}
