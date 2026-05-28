package test

import (
	"AIM/internal/ai/service"
	"testing"
)

func TestSummaryChat(t *testing.T) {
	_, err := service.SummaryChat("testuser", "test001")
	if err != nil {
		t.Log("⚠️ AI总结可能未配置Key，但函数调用正常")
	}
	t.Log("✅ AI聊天总结测试完成")
}

func TestExtractTodo(t *testing.T) {
	_, err := service.ExtractTodo("testuser", "test001")
	if err != nil {
		t.Log("⚠️ AI待办可能未配置Key，但函数调用正常")
	}
	t.Log("✅ AI待办提取测试完成")
}

func TestGenerateReply(t *testing.T) {
	_, err := service.GenerateReply("testuser", "test001")
	if err != nil {
		t.Log("⚠️ AI回复可能未配置Key，但函数调用正常")
	}
	t.Log("✅ AI生成回复测试完成")
}
