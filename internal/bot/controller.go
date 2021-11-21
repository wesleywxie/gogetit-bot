package bot

import (
	"fmt"
	"github.com/wesleywxie/gogetit/internal/config"
	"github.com/wesleywxie/gogetit/internal/model"
	"github.com/wesleywxie/gogetit/internal/task"
	"github.com/wesleywxie/gogetit/internal/task/ytb"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v3"
)

func startCmdCtr(c tb.Context) error {
	user, _ := model.FindOrCreateUserByTelegramID(c.Chat().ID)
	zap.S().Infof("/start user_id: %d telegram_id: %d", user.ID, user.TelegramID)
	return c.Send(fmt.Sprintf("你好，欢迎使用%v。", config.ProjectName))
}

func ytbCmdCtr(c tb.Context) (err error) {
	url := GetHyperlinkFromMessage(c.Message())

	zap.S().Debugw("Received ytb download command",
		"url", url,
		)

	gen := make(chan string)
	download := make(chan string)

	msg, err := B.Send(c.Chat(), "正在下载...")

	// generate filename
	go ytb.GetFilename(c, msg, url, gen)

	// execute download and store
	go ytb.ExecDownload(c, msg, url, gen, download)

	// upload with rclone
	go task.Sync(c, msg, download)

	return
}

func helpCmdCtr(c tb.Context) error {
	message := `
命令： 
/ytb yt-dlp 下载 
/help 帮助
/version 查看当前bot版本
详细使用方法请看：https://github.com/wesleywxie/gogetit
`
	return c.Send(message)
}

func versionCmdCtr(c tb.Context) error {
	return c.Send(config.AppVersionInfo())
}
