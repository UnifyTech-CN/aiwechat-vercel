package config

import (
	"os"
	"strings"
)

const (
	Gemini_Welcome_Reply_Key = "geminiWelcomeReply"
	Gemini_Key               = "geminiKey"
	Gemini_Model_Key         = "geminiModel"
	DefaultGeminiWelcome     = "我是gemini，开始聊天吧！"
	DefaultGeminiModel       = "gemini-1.5-flash"
)

// GetGeminiWelcomeReply returns the welcome message for Gemini bot
func GetGeminiWelcomeReply() string {
	if reply := os.Getenv(Gemini_Welcome_Reply_Key); reply != "" {
		return strings.TrimSpace(reply)
	}
	return DefaultGeminiWelcome
}

// GetGeminiKey returns the Gemini API key
func GetGeminiKey() string {
	return strings.TrimSpace(os.Getenv(Gemini_Key))
}

// GetGeminiModel returns the Gemini model name
func GetGeminiModel() string {
	if model := os.Getenv(Gemini_Model_Key); model != "" {
		return strings.TrimSpace(model)
	}
	return DefaultGeminiModel
}

// IsGeminiConfigured checks if Gemini is properly configured
func IsGeminiConfigured() bool {
	return GetGeminiKey() != ""
}