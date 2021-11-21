package bot

import (
	"fmt"
	"github.com/wesleywxie/gogetit/internal/config"
	"github.com/wesleywxie/gogetit/internal/model"
	"github.com/wesleywxie/gogetit/internal/util"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v3"
	"os"
	"os/exec"
)

func startCmdCtr(c tb.Context) error {
	user, _ := model.FindOrCreateUserByTelegramID(c.Chat().ID)
	zap.S().Infof("/start user_id: %d telegram_id: %d", user.ID, user.TelegramID)
	return c.Send(fmt.Sprintf("你好，欢迎使用%v。", config.ProjectName))
}

func ytbCmdCtr(c tb.Context) error {
	url := GetHyperlinkFromMessage(c.Message())

	zap.S().Debugw("Received ytb download command",
		"url", url)

	args := util.BuildYtdlpArgs(url)

	cmd := exec.Command("yt-dlp", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		zap.S().Warnw("Failed to download",
			"url", url,
			"error", err.Error(),
			)
		return c.Send("下载失败")
	}

	return c.Send("下载完成")
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
