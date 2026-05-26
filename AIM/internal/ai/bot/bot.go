package bot

import (
	"AIM/internal/ai/adapter"
	"AIM/internal/ai/billing"
	"AIM/internal/comet/protocol"
	"AIM/internal/comet/push"
	"AIM/internal/storage/repo"
	"errors"
	"fmt"
	"reflect"
)

// GetClient 获取AI客户端
func GetClient(uid string) (adapter.AIClient, error) {
	cfg, err := billing.GetUserConfig(uid)
	if err != nil {
		return adapter.NewOpenAIClient("PLATFORM_KEY", "gpt-3.5-turbo"), nil
	}
	switch cfg.Platform {
	case "openai":
		return adapter.NewOpenAIClient(cfg.APIKey, "gpt-3.5-turbo"), nil
	case "doubao":
		return adapter.NewDoubaoClient(cfg.APIKey), nil
	}
	return nil, errors.New("unsupported platform")
}

// SendReply 发送消息并获取回复
func SendReply(uid, content string) (string, error) {
	fmt.Println("DEBUG: SendReply start, uid:", uid, "content:", content)

	ok, err := billing.CheckLimit(uid)
	if err != nil || !ok {
		fmt.Println("DEBUG: CheckLimit failed, err:", err, "ok:", ok)
		return "", errors.New("quota exceeded")
	}
	fmt.Println("DEBUG: CheckLimit passed")

	client, err := GetClient(uid)
	if err != nil {
		fmt.Println("DEBUG: GetClient failed, err:", err)
		return "", err
	}
	fmt.Println("DEBUG: GetClient success, client type:", reflect.TypeOf(client))

	reply, err := client.Chat(content)
	if err != nil {
		fmt.Println("DEBUG: client.Chat failed, err:", err)
		return "", err
	}
	fmt.Println("DEBUG: client.Chat success, reply len:", len(reply), "reply:", reply)

	_ = billing.AddUsage(uid, int64(len(content)/3))
	return reply, nil
}

// SendAIReply 私聊AI回复
func SendAIReply(uid string, content string) {
	reply, err := SendReply(uid, content)
	if err != nil {
		return
	}

	push.PushToUser(uid, &protocol.Message{
		Op:      protocol.OpPrivateChat,
		FromUID: "bot",
		ToUID:   uid,
		Content: reply,
	})
}

// SendGroupAIReply 群聊AI回复
func SendGroupAIReply(groupID int64, content string) {
	reply, err := SendReply("bot", content)
	if err != nil {
		return
	}

	// 获取群成员ID列表
	var uids []string
	repo.GetGroupMemberUids(groupID, &uids)

	// 推送群消息
	push.PushToGroup(uids, &protocol.Message{
		Op:      protocol.OpGroupChat,
		FromUID: "bot",
		GroupID: groupID,
		Content: reply,
	})
}
