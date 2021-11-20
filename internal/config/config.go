package config

import (
	"fmt"
	"github.com/spf13/viper"
	tb "gopkg.in/tucnak/telebot.v3"
)


type RunType string

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	ProjectName          string = "gogetit"
	Socks5               string

	// SQLitePath relative path to SQLite db file
	SQLitePath           string


	// BotToken telegram bot token
	BotToken             string

	// TelegramEndpoint telegram api endpoint, empty by default
	TelegramEndpoint 	 string = tb.DefaultApiURL

	RunMode RunType = ReleaseMode

	// AllowUsers 允许使用bot的用户
	AllowUsers []int64

	// DBLogMode 是否打印数据库日志
	DBLogMode bool = false
)

const (
	TestMode    RunType = "Test"
	ReleaseMode RunType = "Release"
)

func AppVersionInfo() (s string) {
	s = fmt.Sprintf("version %v, commit %v, built at %v", version, commit, date)
	return
}

// GetString get string config value by key
func GetString(key string) string {
	var value string
	if viper.IsSet(key) {
		value = viper.GetString(key)
	}

	return value
}
