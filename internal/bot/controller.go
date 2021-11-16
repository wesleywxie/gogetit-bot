package bot

import (
	"github.com/wesleywxie/gogetit/internal/config"
	tb "gopkg.in/tucnak/telebot.v3"
)

func helpCmdCtr(c tb.Context) error {
	message := `
命令：
/help 帮助
/version 查看当前bot版本
详细使用方法请看：https://github.com/wesleywxie/gogetit
`
	return c.Send(message)
}

func versionCmdCtr(c tb.Context) error {
	return c.Send(config.AppVersionInfo())
}
