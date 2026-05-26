package api

import (
	"AIM/internal/ai/billing"
	"AIM/internal/ai/bot"
	"AIM/internal/ai/service"
	"AIM/internal/comet/protocol"
	"AIM/internal/comet/push"
	"AIM/internal/common/utils"
	"AIM/internal/logic/chat/message"
	"AIM/internal/logic/relation/friend"
	"fmt"
	"net/http"
	"strconv"

	"AIM/internal/common/error"
	"AIM/internal/logic/relation/group"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler() *Handler {
	return &Handler{
		service: NewService(),
	}
}

// ======================
// 用户注册
// ======================
func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	user, err := h.service.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": user})
}

// ======================
// 用户登录
// ======================
func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	user, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	// 调用你写好的 GenerateToken 生成 Token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "生成token失败"})
		return
	}

	// 返回 token + 用户信息
	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "登录成功",
		"token": token, // 这里就是你要的 token
		"user":  user,
	})
}

// ======================
// 发送消息（旧接口兼容）
// ======================
func (h *Handler) SendMessage(c *gin.Context) {
	var req message.MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	err := message.SendPrivateMsg(req.FromUID, req.ToUID, req.Content, protocol.MsgTypeText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "发送失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "发送成功"})
}

// ======================
// 获取聊天记录
// ======================
func (h *Handler) GetChatHistory(c *gin.Context) {
	uid1 := c.Query("uid1")
	uid2 := c.Query("uid2")
	limit := 50

	list, err := message.List(uid1, uid2, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": list})
}

// ======================
// 发送私聊消息（新版）
// ======================
func SendPrivateMessage(c *gin.Context) {
	fromUID, exists := c.Get("userID")
	if !exists {
		error.Fail(c, 1003, "登录已经过期了")
		return
	}

	type Req struct {
		ToUID   string `json:"to_uid" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}

	err := message.SendPrivateMsg(fromUID.(string), req.ToUID, req.Content, protocol.MsgTypeText)
	if err != nil {
		error.Fail(c, 2001, "消息发送失败")
		return
	}

	error.Success(c, nil)
}

// ======================
// 群聊相关
// ======================
type CreateGroupReq struct {
	GroupName string `json:"group_name" binding:"required"`
	Desc      string `json:"desc"`
}

type JoinGroupReq struct {
	GroupID string `json:"group_id" binding:"required"`
}

type SendGroupMsgReq struct {
	GroupID int64  `json:"group_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// 创建群聊
func CreateGroup(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已经过期了")
		return
	}
	var req CreateGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	gInfo, err := group.CreateGroup(req.GroupName, uid.(string), req.Desc)
	if err != nil {
		error.Fail(c, 6001, "创建群聊失败")
		return
	}
	error.Success(c, gInfo)
}

// 加入群聊
func JoinGroup(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已经过期了")
		return
	}
	var req JoinGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	groupID, err := strconv.ParseInt(req.GroupID, 10, 64)
	if err != nil {
		error.Fail(c, 1004, "group_id 格式错误，必须为数字")
		return
	}
	err = group.JoinGroup(groupID, uid.(string))
	if err != nil {
		error.Fail(c, 6002, "加入群聊失败")
		return
	}
	error.Success(c, nil)
}

// 获取我的群列表
func GetMyGroupList(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已经过期了")
		return
	}
	list, err := group.ListMyGroup(uid.(string))
	if err != nil {
		error.Fail(c, 6003, "获取群列表失败")
		return
	}
	error.Success(c, list)
}

// 发送群消息
func SendGroupMsg(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已经过期了")
		return
	}
	var req SendGroupMsgReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	err := group.SendGroupMsg(uid.(string), req.GroupID, req.Content)
	if err != nil {
		error.Fail(c, 2001, "消息发送失败")
		return
	}
	error.Success(c, nil)
}

func ChatGroupAI(c *gin.Context) {

	var req AIGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}

	go bot.SendGroupAIReply(req.GroupID, req.Content)
	error.Success(c, "AI群回复中")
}

type AIReq struct {
	Content string `json:"content" binding:"required"`
}

type AIGroupReq struct {
	GroupID int64  `json:"group_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// 私聊调用AI
func ChatAI(c *gin.Context) {
	var req AIReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}

	reply, err := bot.SendReply(uid.(string), req.Content)
	if err != nil {
		error.Fail(c, 500, "AI回复失败: "+err.Error())
		return
	}
	fmt.Println("DEBUG reply len:", len(reply), "reply:", reply)

	// 异步推送，不阻塞、不返回
	go func() {
		push.PushToUser(uid.(string), &protocol.Message{
			Op:      protocol.OpPrivateChat,
			FromUID: "bot",
			ToUID:   uid.(string),
			Content: reply,
		})
	}()

	// 只在这里返回一次！
	error.Success(c, reply)
}

// 群聊调用AI

// ======================
// 好友功能
// ======================
type FriendReq struct {
	FriendID string `json:"friend_id" binding:"required"`
	Remark   string `json:"remark"`
}

// 添加好友
func AddFriend(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已经过期了")
		return
	}

	var req FriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}

	err := friend.AddFriend(uid.(string), req.FriendID, req.Remark)
	if err != nil {
		error.Fail(c, 4001, "添加好友失败")
		return
	}
	error.Success(c, nil)
}

// 好友列表
func ListFriend(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已经过期了")
		return
	}

	list, err := friend.ListFriend(uid.(string))
	if err != nil {
		error.Fail(c, 4002, "获取好友列表失败")
		return
	}
	error.Success(c, list)
}

// 删除好友
func DelFriend(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已经过期了")
		return
	}

	var req FriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}

	err := friend.DelFriend(uid.(string), req.FriendID)
	if err != nil {
		error.Fail(c, 4003, "删除好友失败")
		return
	}
	error.Success(c, nil)
}

// ======================
// 消息已读回执
// ======================
type ReadReceiptReq struct {
	FromUID string `json:"from_uid" binding:"required"`
}

func ReadMessage(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已经过期了")
		return
	}

	var req ReadReceiptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}

	err := message.ReadMessages(req.FromUID, uid.(string))
	if err != nil {
		error.Fail(c, 2002, "标记已读失败")
		return
	}

	error.Success(c, nil)
}

// ======================
// 对方正在输入
// ======================
type TypingReq struct {
	TargetUID string `json:"target_uid" binding:"required"`
}

func SendTyping(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已经过期了")
		return
	}

	var req TypingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}

	err := message.SendTypingStatus(uid.(string), req.TargetUID)
	if err != nil {
		error.Fail(c, 2003, "状态发送失败")
		return
	}
	error.Success(c, nil)
}

// ======================
// 多类型消息（文本 / 图片 / 文件 / 语音）
// ======================
type SendTextReq struct {
	ToUID   string `json:"to_uid" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type SendImageReq struct {
	ToUID    string `json:"to_uid" binding:"required"`
	ImageUrl string `json:"image_url" binding:"required"`
}

type SendFileReq struct {
	ToUID   string `json:"to_uid" binding:"required"`
	FileUrl string `json:"file_url" binding:"required"`
}

type SendVoiceReq struct {
	ToUID    string `json:"to_uid" binding:"required"`
	VoiceUrl string `json:"voice_url" binding:"required"`
}

// 发送文本
func SendText(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}

	var req SendTextReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}

	err := message.SendPrivateMsg(uid.(string), req.ToUID, req.Content, protocol.MsgTypeText)
	if err != nil {
		error.Fail(c, 2001, "发送失败")
		return
	}
	error.Success(c, nil)
}

// 发送图片
func SendImage(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}

	var req SendImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}

	err := message.SendPrivateMsg(uid.(string), req.ToUID, req.ImageUrl, protocol.MsgTypeImage)
	if err != nil {
		error.Fail(c, 2001, "发送失败")
		return
	}
	error.Success(c, nil)
}

// 发送文件
func SendFile(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}

	var req SendFileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}

	err := message.SendPrivateMsg(uid.(string), req.ToUID, req.FileUrl, protocol.MsgTypeFile)
	if err != nil {
		error.Fail(c, 2001, "发送失败")
		return
	}
	error.Success(c, nil)
}

// 发送语音
func SendVoice(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}

	var req SendVoiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}

	err := message.SendPrivateMsg(uid.(string), req.ToUID, req.VoiceUrl, protocol.MsgTypeVoice)
	if err != nil {
		error.Fail(c, 2001, "发送失败")
		return
	}
	error.Success(c, nil)
}
func GetHistoryRoam(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}

	friendUID := c.Query("friend_uid")
	if friendUID == "" {
		error.Fail(c, 1004, "参数错误")
		return
	}

	list, err := message.List(uid.(string), friendUID, 100)
	if err != nil {
		error.Fail(c, 2004, "获取漫游消息失败")
		return
	}
	error.Success(c, list)
}

// 2. 按关键词搜索消息
type SearchKeywordReq struct {
	FriendUID string `json:"friend_uid" binding:"required"`
	Keyword   string `json:"keyword" binding:"required"`
}

func SearchByKeyword(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}

	var req SearchKeywordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}

	list, err := message.SearchKeyword(uid.(string), req.FriendUID, req.Keyword)
	if err != nil {
		error.Fail(c, 2005, "搜索失败")
		return
	}
	error.Success(c, list)
}

// 3. 按时间范围搜索消息
type SearchTimeReq struct {
	FriendUID string `json:"friend_uid" binding:"required"`
	StartTime int64  `json:"start_time" binding:"required"`
	EndTime   int64  `json:"end_time" binding:"required"`
}

func SearchByTime(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}

	var req SearchTimeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}

	list, err := message.SearchTime(uid.(string), req.FriendUID, req.StartTime, req.EndTime)
	if err != nil {
		error.Fail(c, 2005, "搜索失败")
		return
	}
	error.Success(c, list)
}

// ======================
// 群组高级管理（全套）
// ======================
type GroupOperatorReq struct {
	GroupID   int64  `json:"group_id" binding:"required"`
	TargetUID string `json:"target_uid" binding:"required"`
}
type GroupMuteReq struct {
	GroupID   int64  `json:"group_id" binding:"required"`
	TargetUID string `json:"target_uid" binding:"required"`
	Minutes   int    `json:"minutes"`
}
type GroupNoticeReq struct {
	GroupID int64  `json:"group_id" binding:"required"`
	Notice  string `json:"notice" binding:"required"`
}

// 踢人
func GroupKick(c *gin.Context) {
	uid, _ := c.Get("userID")
	var req GroupOperatorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	if err := group.KickMember(uid.(string), req.GroupID, req.TargetUID); err != nil {
		error.Fail(c, 6004, err.Error())
		return
	}
	error.Success(c, nil)
}

// 禁言
func GroupMute(c *gin.Context) {
	uid, _ := c.Get("userID")
	var req GroupMuteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	if err := group.MuteMember(uid.(string), req.GroupID, req.TargetUID, req.Minutes); err != nil {
		error.Fail(c, 6005, err.Error())
		return
	}
	error.Success(c, nil)
}

// 解除禁言
func GroupUnMute(c *gin.Context) {
	uid, _ := c.Get("userID")
	var req GroupOperatorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	if err := group.UnMuteMember(uid.(string), req.GroupID, req.TargetUID); err != nil {
		error.Fail(c, 6006, err.Error())
		return
	}
	error.Success(c, nil)
}

// 转让群主
func GroupTransferOwner(c *gin.Context) {
	uid, _ := c.Get("userID")
	var req GroupOperatorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	if err := group.TransferOwner(uid.(string), req.GroupID, req.TargetUID); err != nil {
		error.Fail(c, 6007, err.Error())
		return
	}
	error.Success(c, nil)
}

// 设置群公告
func GroupSetNotice(c *gin.Context) {
	uid, _ := c.Get("userID")
	var req GroupNoticeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	if err := group.UpdateNotice(uid.(string), req.GroupID, req.Notice); err != nil {
		error.Fail(c, 6008, err.Error())
		return
	}
	error.Success(c, nil)
}

// 获取群公告
func GroupGetNotice(c *gin.Context) {
	gid := c.Query("group_id")
	id, _ := strconv.ParseInt(gid, 10, 64)
	notice, err := group.GetNotice(id)
	if err != nil {
		error.Fail(c, 6009, "获取失败")
		return
	}
	error.Success(c, gin.H{"notice": notice})
}

// 设置管理员
func GroupSetAdmin(c *gin.Context) {
	uid, _ := c.Get("userID")
	var req GroupOperatorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	if err := group.SetAdmin(uid.(string), req.GroupID, req.TargetUID); err != nil {
		error.Fail(c, 6010, err.Error())
		return
	}
	error.Success(c, nil)
}

// 撤销管理员
func GroupRemoveAdmin(c *gin.Context) {
	uid, _ := c.Get("userID")
	var req GroupOperatorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	if err := group.RemoveAdmin(uid.(string), req.GroupID, req.TargetUID); err != nil {
		error.Fail(c, 6011, err.Error())
		return
	}
	error.Success(c, nil)
}

// ====================== 好友分组 ======================
type FriendGroupReq struct {
	GroupName string `json:"group_name" binding:"required"`
}
type FriendMoveGroupReq struct {
	FriendID string `json:"friend_id" binding:"required"`
	GroupID  int64  `json:"group_id" binding:"required"`
}
type FriendGroupIDReq struct {
	GroupID int64 `json:"group_id" binding:"required"`
}

// 创建好友分组
func CreateFriendGroup(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}
	var req FriendGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	if err := friend.CreateGroup(uid.(string), req.GroupName); err != nil {
		error.Fail(c, 4005, "创建分组失败")
		return
	}
	error.Success(c, nil)
}

// 获取我的分组列表
func ListFriendGroup(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}
	list, err := friend.ListGroup(uid.(string))
	if err != nil {
		error.Fail(c, 4006, "获取分组失败")
		return
	}
	error.Success(c, list)
}

// 删除分组
func DeleteFriendGroup(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}
	var req FriendGroupIDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	if err := friend.DeleteGroup(uid.(string), req.GroupID); err != nil {
		error.Fail(c, 4007, "删除分组失败")
		return
	}
	error.Success(c, nil)
}

// 移动好友到分组
func MoveFriendToGroup(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}
	var req FriendMoveGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	if err := friend.MoveFriend(uid.(string), req.FriendID, req.GroupID); err != nil {
		error.Fail(c, 4008, "移动失败")
		return
	}
	error.Success(c, nil)
}

// 按分组获取好友
func ListFriendByGroup(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}
	var req FriendGroupIDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	list, err := friend.ListByGroup(uid.(string), req.GroupID)
	if err != nil {
		error.Fail(c, 4009, "获取失败")
		return
	}
	error.Success(c, list)
}

// ====================== AI 配置 ======================
type SetAIConfigReq struct {
	Platform string `json:"platform" binding:"required"`
	APIKey   string `json:"api_key" binding:"required"`
}

func SetAIConfig(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}
	var req SetAIConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	if err := billing.SaveUserConfig(uid.(string), req.Platform, req.APIKey); err != nil {
		error.Fail(c, 5001, "保存配置失败")
		return
	}
	error.Success(c, nil)
}

// ====================== AI 工具功能 ======================
type ChatReq struct {
	Content string `json:"content" binding:"required"`
}
type ChatTargetReq struct {
	FriendUID string `json:"friend_uid" binding:"required"`
}

func AIChat(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}
	var req ChatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	reply, err := service.Chat(uid.(string), req.Content)
	if err != nil {
		error.Fail(c, 5002, "AI 调用失败")
		return
	}
	error.Success(c, reply)
}

func AISummary(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}
	var req ChatTargetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	summary, err := service.SummaryChat(uid.(string), req.FriendUID)
	if err != nil {
		error.Fail(c, 5003, "生成总结失败")
		return
	}
	error.Success(c, summary)
}

func AIExtractTodo(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}
	var req ChatTargetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	todo, err := service.ExtractTodo(uid.(string), req.FriendUID)
	if err != nil {
		error.Fail(c, 5004, "提取待办失败")
		return
	}
	error.Success(c, todo)
}

func AIGenerateReply(c *gin.Context) {
	uid, ok := c.Get("userID")
	if !ok {
		error.Fail(c, 1003, "登录已过期")
		return
	}
	var req ChatTargetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		error.Fail(c, 1004, "参数错误")
		return
	}
	replies, err := service.GenerateReply(uid.(string), req.FriendUID)
	if err != nil {
		error.Fail(c, 5005, "生成回复失败")
		return
	}
	error.Success(c, replies)
}
