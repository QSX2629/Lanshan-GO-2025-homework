package main

import (
	"AIM/internal/api"
	"AIM/internal/comet/connection"
	"AIM/internal/common/config"
	"AIM/internal/common/logger"
	"AIM/internal/common/middleware"
	"AIM/internal/storage/mysql"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	logger.Init()
	mysql.Init()

	r := gin.Default()
	h := api.NewHandler()

	// --------------------------
	// 公共路由（不需要登录）
	// --------------------------
	r.POST("/user/register", h.Register)
	r.POST("/user/login", h.Login)
	r.GET("/chat/history", h.GetChatHistory)

	// --------------------------
	// WebSocket 长连接（实时聊天）
	// --------------------------
	r.GET("/ws", func(c *gin.Context) {
		connection.HandleConnection(c.Writer, c.Request)
	})

	// --------------------------
	// 需要登录鉴权的路由
	// --------------------------
	authGroup := r.Group("/")
	authGroup.Use(middleware.Auth())
	{
		// 新版发送消息：存库 + 实时推送
		authGroup.POST("/chat/send", api.SendPrivateMessage)
		authGroup.POST("/group/create", api.CreateGroup) // 创建群
		authGroup.POST("/group/join", api.JoinGroup)     // 加入群
		authGroup.GET("/group/list", api.GetMyGroupList) // 我的群列表
		authGroup.POST("/group/send", api.SendGroupMsg)  // 发送群消息
		authGroup.POST("/ai/chat", api.ChatAI)           // 私聊AI
		authGroup.POST("/ai/group", api.ChatGroupAI)
		// 已读回执
		authGroup.POST("/chat/read", api.ReadMessage)
		authGroup.POST("/chat/typing", api.SendTyping)
		authGroup.POST("/chat/send/text", api.SendText)
		authGroup.POST("/chat/send/image", api.SendImage)
		authGroup.POST("/chat/send/file", api.SendFile)
		authGroup.POST("/chat/send/voice", api.SendVoice)
		authGroup.GET("/chat/roam", api.GetHistoryRoam)
		authGroup.POST("/chat/search/keyword", api.SearchByKeyword)
		authGroup.POST("/chat/search/time", api.SearchByTime)
		// 群组高级管理
		authGroup.POST("/group/kick", api.GroupKick)
		authGroup.POST("/group/mute", api.GroupMute)
		authGroup.POST("/group/unmute", api.GroupUnMute)
		authGroup.POST("/group/transfer", api.GroupTransferOwner)
		authGroup.POST("/group/notice/set", api.GroupSetNotice)
		authGroup.GET("/group/notice/get", api.GroupGetNotice)
		authGroup.POST("/group/admin/set", api.GroupSetAdmin)
		authGroup.POST("/group/admin/remove", api.GroupRemoveAdmin) // 群聊AI
		// 好友分组
		authGroup.POST("/friend/group/create", api.CreateFriendGroup)
		authGroup.GET("/friend/group/list", api.ListFriendGroup)
		authGroup.POST("/friend/group/delete", api.DeleteFriendGroup)
		authGroup.POST("/friend/group/move", api.MoveFriendToGroup)
		authGroup.POST("/friend/group/friends", api.ListFriendByGroup)
		// AI 相关
		authGroup.POST("/ai/config/set", api.SetAIConfig)

		authGroup.POST("/ai/summary", api.AISummary)
		authGroup.POST("/ai/todo", api.AIExtractTodo)
		authGroup.POST("/ai/reply/generate", api.AIGenerateReply)

	}

	// 启动服务
	r.Run(":8080")
}
