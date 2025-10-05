package api

import (
	"fmt"
	"net/http"

	"github.com/pwh-pwh/aiwechat-vercel/chat"
	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func Wx(rw http.ResponseWriter, req *http.Request) {
	// 打印请求信息用于调试
	fmt.Printf("收到请求: %s %s\n", req.Method, req.URL.Path)
	
	// 验证环境变量是否设置
	token := config.GetWxToken()
	appID := config.GetWxAppId()
	aesKey := config.GetWxEncodingAESKey()
	
	if token == "" {
		fmt.Println("错误: WX_TOKEN 未设置")
	} else {
		fmt.Printf("WX_TOKEN 已设置: %s\n", token)
	}
	
	if appID == "" {
		fmt.Println("错误: WX_APP_ID 未设置")
	} else {
		fmt.Printf("WX_APP_ID 已设置\n")
	}
	
	if aesKey == "" {
		fmt.Println("错误: WX_ENCODING_AES_KEY 未设置")
	} else {
		fmt.Printf("WX_ENCODING_AES_KEY 长度: %d\n", len(aesKey))
	}
	
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &offConfig.Config{
		AppID:          appID,
		AppSecret:      config.GetWxAppSecret(),
		Token:          token,
		EncodingAESKey: aesKey,
		Cache:          memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)

	// 传入request和responseWriter
	server := officialAccount.GetServer(req, rw)
	
	// 对于GET请求（微信验证请求），确保验证开启
	if req.Method == "GET" {
		fmt.Println("处理微信验证请求 (GET)")
		// SkipValidate设置为false表示进行签名验证
		server.SkipValidate(false)
	} else {
		// 对于POST请求（消息推送），根据安全模式决定是否验证
		fmt.Println("处理微信消息推送 (POST)")
		server.SkipValidate(false)
	}
	
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		//回复消息：演示回复用户发送的消息
		replyMsg := handleWxMessage(msg)
		text := message.NewText(replyMsg)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Printf("处理请求时出错: %v\n", err)
		return
	}
	//发送回复的消息
	server.Send()
	fmt.Println("请求处理完成")
}

func handleWxMessage(msg *message.MixMessage) (replyMsg string) {
	msgType := msg.MsgType
	msgContent := msg.Content
	userId := string(msg.FromUserName)

	// Check if user is authenticated (only if ADDME_PASSWORD is set)
	if config.GetAddMePassword() != "" && !config.IsUserAuthenticated(userId) {
		if msgType == message.MsgTypeText {
			// Only allow /addme command for non-authenticated users
			if msgContent == "/addme" || len(msgContent) > len("/addme") && msgContent[:len("/addme")] == "/addme" {
				bot := chat.GetChatBot(config.GetUserBotType(userId))
				replyMsg = bot.Chat(userId, msgContent)
			} else {
				replyMsg = "功能还在开发中"
			}
		} else {
			replyMsg = "功能还在开发中"
		}
		return
	}

	bot := chat.GetChatBot(config.GetUserBotType(userId))
	if msgType == message.MsgTypeText {
		replyMsg = bot.Chat(userId, msgContent)
	} else {
		replyMsg = bot.HandleMediaMsg(msg)
	}

	return
}
