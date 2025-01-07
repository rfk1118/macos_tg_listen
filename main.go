package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime/debug"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("\n🚀 启动 Telegram 机器人服务\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// 加载配置文件
	config, err := LoadConfig("config.yaml")
	if err != nil {
		log.Printf("\n❌ 加载配置文件失败\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"+
			"🔴 错误信息: %v\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n",
			err)
		time.Sleep(5 * time.Second)
		main()
		return
	}

	// 添加全局的 panic 恢复
	defer func() {
		if r := recover(); r != nil {
			log.Printf("\n❌ 程序发生严重错误\n"+
				"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"+
				"🔴 错误信息: %v\n"+
				"📑 堆栈信息:\n%s\n"+
				"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n",
				r, debug.Stack())
			log.Println("⏳ 5秒后尝试重启服务...")
			time.Sleep(5 * time.Second)
			main()
		}
	}()

	// 使用配置创建 bot
	bot, err := tgbotapi.NewBotAPI(config.Telegram.BotToken)
	if err != nil {
		log.Printf("\n❌ 初始化机器人失败\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"+
			"🔴 错误信息: %v\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n",
			err)
		time.Sleep(5 * time.Second)
		main()
		return
	}

	bot.Debug = config.Telegram.Debug
	log.Printf("\n✅ 机器人初始化成功\n"+
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"+
		"👾 机器人名称: %s\n"+
		"🆔 机器人 ID: %d\n"+
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n",
		bot.Self.UserName, bot.Self.ID)

	fmt.Println("接收机器人已启动...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = config.Telegram.Timeout

	updates := bot.GetUpdatesChan(u)

	// 主消息处理循环
	for update := range updates {
		// 为每个消息处理添加 recover
		go func(update tgbotapi.Update) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("\n⚠️ 处理消息时发生错误\n"+
						"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"+
						"🔴 错误信息: %v\n"+
						"📑 堆栈信息:\n%s\n"+
						"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n",
						r, debug.Stack())
				}
			}()

			if update.Message != nil {
				handleMessage(update.Message)
			}
			if update.ChannelPost != nil {
				handleChannelPost(update.ChannelPost)
			}
		}(update)
	}
}

func handleMessage(message *tgbotapi.Message) {
	if message == nil {
		log.Println("收到空消息")
		return
	}

	chatID := message.Chat.ID
	var senderName string

	// 安全地获取发送者信息
	if message.Chat.IsPrivate() {
		if message.From != nil {
			senderName = message.From.UserName
			if senderName == "" {
				senderName = fmt.Sprintf("%s %s", message.From.FirstName, message.From.LastName)
			}
		} else {
			senderName = "未知用户"
		}
	} else {
		if message.SenderChat != nil {
			senderName = message.SenderChat.Title
		} else if message.Chat != nil {
			senderName = message.Chat.Title
		} else {
			senderName = "未知来源"
		}
	}

	log.Printf("\n📩 新消息\n"+
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"+
		"👤 发送者: %s\n"+
		"🆔 聊天 ID: %d\n"+
		"💬 内容: %s\n"+
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n",
		senderName, chatID, message.Text)
	say(message.Text)
}

func handleChannelPost(post *tgbotapi.Message) {
	if post == nil {
		log.Println("收到空的频道消息")
		return
	}

	chatID := post.Chat.ID
	channelTitle := "未知频道"
	if post.Chat != nil {
		channelTitle = post.Chat.Title
	}

	log.Printf("\n📢 频道消息\n"+
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"+
		"📺 频道名称: %s\n"+
		"🆔 频道 ID: %d\n"+
		"💬 消息内容: %s\n"+
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n",
		channelTitle, chatID, post.Text)

	say(post.Text)
}

func say(text string) {
	if text == "" {
		log.Println("🔇 语音播报: 收到空消息，跳过播报")
		return
	}

	cmd := exec.Command("say", text)
	if err := cmd.Run(); err != nil {
		log.Printf("\n❌ 语音播报失败\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"+
			"🔴 错误信息: %v\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n",
			err)
	}
}
