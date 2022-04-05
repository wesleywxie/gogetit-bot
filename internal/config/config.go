package config

import (
	"fmt"
	"github.com/spf13/viper"
	tb "gopkg.in/telebot.v3"
)

type RunType string

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	ProjectName string = "gogetit"
	Socks5      string
	// SQLitePath relative path to SQLite db file
	SQLitePath string
	// BotToken telegram bot token
	BotToken string
	// TelegramEndpoint telegram api endpoint, empty by default
	TelegramEndpoint string = tb.DefaultApiURL
	// UpdateInterval 刷新间隔
	UpdateInterval int = 10

	// UserAgent as the user agent for downloading task
	UserAgent string
	// OutputDir folder to store downloaded files
	OutputDir string
	// CookieFile cookie file exported for specific website
	CookieFile string
	// MaxThreadNum max thread count for aria2c
	MaxThreadNum int

	// AutoUpload when download finished
	AutoUpload bool = false
	// AutoUploadDrive driver name for rclone to upload to
	AutoUploadDrive string
	// AutoUploadDir dir for rclone to upload to
	AutoUploadDir string = ProjectName

	RunMode = ReleaseMode

	// AllowUsers 允许使用bot的用户
	AllowUsers []int64

	// DBLogMode 是否打印数据库日志
	DBLogMode = false
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
