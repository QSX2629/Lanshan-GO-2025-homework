package protocol

// Message 前端 <-> 后端 统一消息协议
type Message struct {
	Op      int    `json:"op"`       // 1=私聊 2=群聊 3=心跳
	FromUID string `json:"from_uid"` // 发送者ID
	ToUID   string `json:"to_uid"`   // 接收者ID
	Content string `json:"content"`  // 消息内容
	GroupID int64  `json:"group_id"`
	MsgType int    `json:"msg_type"`
}

// 操作类型常量
const (
	OpPrivateChat  = 1 // 私聊消息
	OpGroupChat    = 2 // 群聊消息
	OpAIChat       = 3 // AI 私聊回复
	OpGroupAIChat  = 4 // AI
	OpUserStatus   = 5
	OpReadReceipt  = 6
	OpTypingStatus = 7

	// 消息类型
	MsgTypeText  = 1 // 文本
	MsgTypeImage = 2 // 图片
	MsgTypeFile  = 3 // 文件
	MsgTypeVoice = 4 // 语音
)
