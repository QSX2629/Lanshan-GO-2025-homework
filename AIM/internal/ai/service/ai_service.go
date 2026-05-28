package service

import (
	"AIM/internal/ai/bot"
	"AIM/internal/ai/context"
	"AIM/internal/ai/prompt"
	"AIM/internal/logic/chat/message"
	"fmt"
)

// Chat 普通对话
func Chat(uid, content string) (string, error) {
	context.AppendContext(uid, "user: "+content)
	reply, err := bot.SendReply(uid, content)
	if err != nil {
		return "", err
	}
	context.AppendContext(uid, "assistant: "+reply)
	return reply, nil
}

// SummaryChat 聊天记录总结
func SummaryChat(uid, friendUID string) (string, error) {
	msgs, err := message.List(uid, friendUID, 50)
	if err != nil {
		return "", err
	}
	var content string
	for _, m := range msgs {
		content += fmt.Sprintf("%s: %s\n", m.FromUID, m.Content)
	}
	p := fmt.Sprintf(prompt.SummaryPrompt, content)
	return bot.SendReply(uid, p)
}

// ExtractTodo 提取待办
func ExtractTodo(uid, friendUID string) (string, error) {
	msgs, err := message.List(uid, friendUID, 50)
	if err != nil {
		return "", err
	}
	var content string
	for _, m := range msgs {
		content += fmt.Sprintf("%s: %s\n", m.FromUID, m.Content)
	}
	p := fmt.Sprintf(prompt.TodoPrompt, content)
	return bot.SendReply(uid, p)
}

// GenerateReply 生成回复候选
func GenerateReply(uid, friendUID string) (string, error) {
	msgs, err := message.List(uid, friendUID, 20)
	if err != nil {
		return "", err
	}
	var content string
	for _, m := range msgs {
		content += fmt.Sprintf("%s: %s\n", m.FromUID, m.Content)
	}
	p := fmt.Sprintf(prompt.ReplyPrompt, content)
	return bot.SendReply(uid, p)
}
