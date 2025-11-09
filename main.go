package main

import (
	"LOL-Assistant/disocrd"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// main 函数是程序的入口点，用于初始化和运行 Discord 机器人。
// 从 .env 文件加载环境变量，初始化 Gemini 客户端后
// 启动 Discord 机器人。机器人运行后，等待 CTRL-C 等中断信号
// 来终止程序。
func main() {
	// 从 .env 文件加载环境变量
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal(err)
	}

	// 初始化 Gemini 客户端
	disocrd.Initialize()

	// 设置 Discord 机器人令牌
	token := fmt.Sprintf("Bot %s", os.Getenv("BOT_TOKEN"))
	// 创建 Discord 会话
	dg, err := discordgo.New(token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// 注册消息处理器
	dg.AddHandler(disocrd.Message)
	// 设置公会消息意图
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// 开始 Discord 连接
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// 输出机器人运行中消息
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// 创建用于等待退出信号的通道
	sc := make(chan os.Signal, 1)
	// 设置退出信号检测
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	// 阻塞直到接收到信号
	<-sc

	// 关闭 Discord 连接
	defer func() {
		err = dg.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
}
