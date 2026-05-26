package test

import (
	"AIM/internal/api"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateGroup(t *testing.T) {
	r := gin.Default()
	// 直接注册包级函数，不是 h.CreateGroup
	r.POST("/group/create", api.CreateGroup)

	body := strings.NewReader(`{"group_name":"测试群"}`)
	req, _ := http.NewRequest("POST", "/group/create", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	t.Log("✅ 创建群聊测试完成")
}

func TestSendGroupMessage(t *testing.T) {
	r := gin.Default()
	r.POST("/group/send", api.SendGroupMsg)

	body := strings.NewReader(`{"group_id":1,"content":"hello"}`)
	req, _ := http.NewRequest("POST", "/group/send", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	t.Log("✅ 发送群消息测试完成")
}
