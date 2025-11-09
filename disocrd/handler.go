package disocrd

import (
	"LOL-Assistant/gemini"
	"LOL-Assistant/league"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// 声明全局变量 Gemini 客户端
var geminiClient gemini.ChatSession

// Initialize 函数在机器人启动时只调用一次
// 用于初始化 Gemini AI 客户端。
func Initialize() {
	geminiClient = gemini.NewGeminiClient()
}

// Message 函数是 Discord 消息事件处理器。
// 接收用户发送的消息并处理，根据条件发送响应。
// 忽略机器人自己发送的消息。
func Message(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 忽略机器人自己的消息
	if m.Author.ID == s.State.User.ID {
		return
	}

	// 处理包含"机器人"的消息
	if strings.Contains(m.Content, "机器人") {
		// 发送"正在生成回答"消息
		delMsg, err := s.ChannelMessageSend(m.ChannelID, "正在生成回答")
		if err != nil {
			log.Println("发送'正在生成回答'失败", err)
			return
		}
		// 使用 Gemini AI 生成用户消息的响应
		resp, err := geminiClient.ChatWithDiscord(context.Background(), m.Content)
		if err != nil {
			// 响应生成失败时发送错误消息
			_, sendErr := s.ChannelMessageSend(m.ChannelID, "无法生成回答")
			if sendErr != nil {
				log.Fatalln("生成回答失败", err)
				return
			}
			log.Println("gemini api 错误", err)
			return
		}

		// 将 Gemini AI 的响应发送到 Discord 频道
		msg, err := s.ChannelMessageSend(m.ChannelID, resp)
		if err != nil {
			// 响应发送失败时发送错误消息
			_, sendErr := s.ChannelMessageSend(m.ChannelID, "无法生成回答")
			if sendErr != nil {
				log.Fatalln("生成回答失败", err)
				return
			}
			log.Println("gemini api 错误", err)
			return
		} else {
			log.Println(msg)
		}

		// 删除之前发送的"正在生成回答"消息
		err = s.ChannelMessageDelete(m.ChannelID, delMsg.ID)
		if err != nil {
			log.Println("消息删除失败", err)
			return
		}
	} else if strings.HasPrefix(m.Content, "分析我最后一场游戏") {
		// 示例: "分析我最后一场游戏|昵称#标签"
		delMsg, err := s.ChannelMessageSend(m.ChannelID, "正在生成回答")
		if err != nil {
			log.Println("发送'正在生成回答'失败", err)
			return
		}

		split := strings.Split(m.Content, "|")
		nickName := strings.Split(split[1], "#")[0]
		tag := strings.Split(split[1], "#")[1]

		gameInfo, puuid, err := league.GetMatch(nickName, tag)
		if err != nil {
			if err != nil {
				// 响应发送失败时发送错误消息
				_, sendErr := s.ChannelMessageSend(m.ChannelID, "无法生成回答")
				if sendErr != nil {
					log.Fatalln("生成回答失败", err)
					return
				}
				log.Println("gemini api 错误", err)
				return
			}
			err = s.ChannelMessageDelete(m.ChannelID, delMsg.ID)
			if err != nil {
				log.Println("消息删除失败", err)
				return
			}
		}
		matchReq := fmt.Sprintf("%s |  puuid: %s, 游戏信息: %s | 使用我的 puuid、昵称和游戏标签来分析我玩的角色信息", m.Content, puuid, gameInfo)

		resp, err := geminiClient.ChatWithDiscord(context.Background(), matchReq)
		if err != nil {
			// 响应生成失败时发送错误消息
			_, sendErr := s.ChannelMessageSend(m.ChannelID, "无法生成回答")
			if sendErr != nil {
				log.Fatalln("生成回答失败", err)
				return
			}
			log.Println("gemini api 错误", err)
			return
		}
		log.Println(resp)

		_, err = s.ChannelMessageSend(m.ChannelID, resp)
		if err != nil {
			_, sendErr := s.ChannelMessageSend(m.ChannelID, "无法生成回答")
			if sendErr != nil {
				log.Fatalln("生成回答失败", err)
				return
			}
			log.Println("gemini api 错误", err)
			return
		}
		err = s.ChannelMessageDelete(m.ChannelID, delMsg.ID)
		if err != nil {
			log.Println("消息删除失败", err)
			return
		}
	}
}
