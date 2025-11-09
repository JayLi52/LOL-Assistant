package gemini

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// ChatSession 结构体用于管理与 Gemini AI 的聊天会话。
type ChatSession struct {
	ChatSession *genai.ChatSession
}

// ChatHistory 是存储聊天记录的全局变量。
var ChatHistory []*genai.Content

// NewGeminiClient 函数创建一个新的 Gemini AI 客户端。
// 加载 API 密钥，初始化模型，配置安全设置，
// 加载 PDF 文件并添加到 ChatHistory 中。
// 最后启动聊天会话并返回。
func NewGeminiClient() ChatSession {
	// 从环境变量加载 API 密钥
	apiKey := os.Getenv("GEMINI_API_KEY")
	ctx := context.Background()

	// 创建 Gemini API 客户端
	var err error
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	// 从环境变量加载系统指令
	instructions := os.Getenv("GEMINI_INSTRUCTIONS")
	// 初始化模型（使用 gemini-2.5-flash-preview-04-17 版本）
	model := client.GenerativeModel("gemini-2.5-flash-preview-04-17")
	// 设置系统指令
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(instructions)},
	}

	// 配置安全设置（所有类别都不进行过滤）
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
	}

	// 读取 pdfs 目录中的文件列表
	files, err := os.ReadDir("pdfs")
	if err != nil {
		log.Fatal(err)
	}

	// 创建文件名数组
	var names []string
	for _, file := range files {
		if !file.IsDir() {
			names = append(names, file.Name())
		}
	}

	// 加载每个 PDF 文件并添加到 ChatHistory
	for _, name := range names {
		wikiData, err := os.ReadFile(fmt.Sprintf("pdfs/%s", name))
		if err != nil {
			log.Fatal(err)
		}
		// 检测文件 MIME 类型
		wikiMimeType := http.DetectContentType(wikiData)

		// 将文件数据以模型角色添加到 ChatHistory
		ChatHistory = append(ChatHistory, &genai.Content{
			Parts: []genai.Part{
				genai.Blob{
					MIMEType: wikiMimeType,
					Data:     wikiData,
				},
			},
			Role: "model",
		})
	}

	// 启动聊天会话
	cs := model.StartChat()

	// 返回 ChatSession 对象
	return ChatSession{
		ChatSession: cs,
	}
}

// ChatWithDiscord 方法将 Discord 接收到的文本发送到 Gemini AI
// 并接收响应返回。对话记录保存在 ChatHistory 中。
// 如果发生错误，返回空字符串和错误。
func (cs ChatSession) ChatWithDiscord(ctx context.Context, text string) (string, error) {
	// 为当前聊天会话设置对话记录
	cs.ChatSession.History = ChatHistory

	// 向 Gemini AI 发送消息
	resp, err := cs.ChatSession.SendMessage(ctx, genai.Text(text))
	if err != nil {
		return "", err
	}

	// 从响应中提取内容
	var content genai.Part
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				content = part
			}
		}
	}

	// 将用户消息添加到对话记录
	ChatHistory = append(ChatHistory, &genai.Content{
		Parts: []genai.Part{
			genai.Text(text),
		},
		Role: "user",
	})

	// 将 AI 响应添加到对话记录
	ChatHistory = append(ChatHistory, &genai.Content{
		Parts: []genai.Part{
			content,
		},
		Role: "model",
	})

	// 以文本格式返回响应
	return string(content.(genai.Text)), nil
}
