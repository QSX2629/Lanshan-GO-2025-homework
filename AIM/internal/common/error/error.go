package error

import "github.com/gin-gonic/gin"

// 错误类型集合：
// 1xxx：用户
// 2xxx：消息
// 3xxx：好友
// 4xxx：群聊
// 5xxx：ai
// 6xxx：系统/存储/服务异常
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

var OK = &Error{
	Code:    0,
	Message: "成功",
}
var (
	UserNotExist     = &Error{Code: 1001, Message: "用户不存在"}
	UserAlreadyExist = &Error{Code: 1102, Message: "用户已经在了"}
	TokenInvalid     = &Error{Code: 1003, Message: "登录已经过期了"}
	ParamError       = &Error{Code: 1004, Message: "参数错误"}
)
var (
	MessageSendFail       = &Error{Code: 2001, Message: "消息发送失败"}
	MessageContentIllegal = &Error{Code: 2002, Message: "非法消息！！！！"}
	MessageNotFound       = &Error{Code: 2003, Message: "消息不存在"}
)
var (
	NotFriend = &Error{Code: 3001, Message: "你们还不是好友"}
)
var (
	GroupNotFound     = &Error{Code: 4001, Message: "没有这个群"}
	NotInGroup        = &Error{Code: 4002, Message: "不在这个群"}
	NoGroupPermission = &Error{Code: 4003, Message: "没有群的权限"}
)
var (
	AIServiceError   = &Error{Code: 5001, Message: "ai服务一场"}
	AIContentExpired = &Error{Code: 5002, Message: "过期了"}
	APINull          = &Error{Code: 5003, Message: "api密钥未设置"}
)
var (
	ServerError    = &Error{Code: 6001, Message: "服务器内部错误"}
	DBError        = &Error{Code: 6002, Message: "数据库异常"}
	RedisError     = &Error{Code: 6003, Message: "缓存异常"}
	ServiceTimeout = &Error{Code: 6004, Message: "服务超时"}
)

type Res struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功返回
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Res{
		Code: 200,
		Msg:  "操作成功",
		Data: data,
	})
}

// Fail 失败返回
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(200, Res{
		Code: code,
		Msg:  msg,
	})
}
