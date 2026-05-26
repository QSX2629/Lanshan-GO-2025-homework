package test

import (
	"AIM/internal/api"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	handler := &api.Handler{} // 初始化你的 Handler
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	return r
}

func TestRegisterAPI(t *testing.T) {
	router := setupRouter()

	// 构造请求体
	body := strings.NewReader(`{"username":"test123","password":"123456"}`)
	req, _ := http.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言响应
	if w.Code != http.StatusOK {
		t.Fatalf("❌ 注册失败，状态码：%d，响应：%s", w.Code, w.Body.String())
	}
	t.Log("✅ 注册接口测试通过")
}

func TestLoginAPI(t *testing.T) {
	router := setupRouter()

	body := strings.NewReader(`{"username":"test123","password":"123456"}`)
	req, _ := http.NewRequest("POST", "/login", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("❌ 登录失败，状态码：%d，响应：%s", w.Code, w.Body.String())
	}
	t.Log("✅ 登录接口测试通过")
}
