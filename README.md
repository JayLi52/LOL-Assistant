# 使用 Go 和 Gemini 构建英雄联盟助手 Discord 机器人

## 安装
```shell
go get -u github.com/bwmarrin/discordgo
go get -u github.com/google/generative-ai-go
go get -u github.com/joho/godotenv
```

## 环境变量
```dotenv
BOT_TOKEN="Discord 机器人令牌"
GEMINI_API_KEY="Gemini API Key"
GEMINI_INSTRUCTIONS="Gemini 指令说明"
RIOT_GAMES_API_KEY="Riot Games API Key"
```
```shell
cp .env.example .env
```
请复制 `.env.example` 文件创建 `.env` 文件。

## 幻灯片
[链接](https://docs.google.com/presentation/d/1Ja5rL4fyNQ3PSRv_hyy23PqtVl0dL-nuZG1BwClwn5k/edit?usp=sharing )

## Riot Games 开放 API
[链接](https://developer.riotgames.com/apis)

## Discord 机器人
[链接](https://discord.com/developers/applications)

## 使用方法
### 运行
```shell
go mod tidy
go run .
```
![run](readme.png)

### 命令
```
机器人 ${提示词}
```
- 可以与 Gemini 进行对话。（参考幻灯片）

```
分析我最后一场游戏|${昵称}#${标签}
```
- 使用 Riot 开放 API 获取最后一场游戏，由 Gemini 进行分析。（参考幻灯片）
