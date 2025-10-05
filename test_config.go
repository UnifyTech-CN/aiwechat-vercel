package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pwh-pwh/aiwechat-vercel/config"
)

func main() {
	// 加载.env文件
	if err := godotenv.Load("conf/.env"); err != nil {
		fmt.Printf("警告: 无法加载.env文件: %v\n", err)
	}
	
	// 运行配置检查
	botType, checkRes := config.CheckAllBotConfig()
	
	// 输出结果
	for bot, status := range checkRes {
		fmt.Printf("%v: %v\n", bot, status)
	}
	fmt.Printf("DEFAULT BOT: %v\n", botType)
}