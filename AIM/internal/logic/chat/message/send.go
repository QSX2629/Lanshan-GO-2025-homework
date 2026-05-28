package message

import (
	"AIM/internal/comet/protocol"
	"AIM/internal/comet/push"
	"AIM/internal/logic/relation/friend"

	"AIM/internal/storage/repo"
	"errors"
)

// SendPrivate 发送单聊消息
func SendPrivateMsg(fromUID, toUID, content string, msgType int) error {
	// 先判断是否好友
	if !friend.IsFriend(fromUID, toUID) {
		return errors.New("你们还不是好友")
	}

	// 存库
	err := repo.CreateMessage(fromUID, toUID, 0, content, msgType)
	if err != nil {
		return err
	}

	// 推送
	push.PushToUser(toUID, &protocol.Message{
		Op:      protocol.OpPrivateChat,
		FromUID: fromUID,
		ToUID:   toUID,
		Content: content,
		MsgType: msgType,
	})
	return nil
}

// ReadMessages 标记消息已读
func ReadMessages(fromUID, toUID string) error {
	return repo.SetMessageRead(fromUID, toUID)
}

// SendTypingStatus 发送正在输入状态
func SendTypingStatus(fromUID, toUID string) error {
	return repo.PushTypingNotify(fromUID, toUID)
}
