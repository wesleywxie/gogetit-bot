package bot

import (
	"bytes"
	"fmt"
	"github.com/wesleywxie/gogetit/internal/config"
	"github.com/wesleywxie/gogetit/internal/model"
	"github.com/wesleywxie/gogetit/internal/util"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v3"
	"io"
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
		"url", url,
		)

	args := util.GetYtdlpFilename(url)
	cmd := exec.Command("yt-dlp", args...)
	out, err := cmd.CombinedOutput()
	zap.S().Debug(string(out))
	if err != nil {
		zap.S().Warnw("Failed to extract filename",
			"url", url,
			"error", err.Error(),
		)
		return c.Send("下载失败")
	}
	filename := string(out)

	args = util.BuildYtdlpArgs(url, filename)

	zap.S().Debugf("Executing command yt-dlp %v", args)
	cmd = exec.Command("yt-dlp", args...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err = cmd.Run()
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())

	if err != nil {
		zap.S().Debugf("Finished command with error output\n %v", errStr)
		zap.S().Warnw("Failed to download",
			"url", url,
			"error", err.Error(),
			)
		return c.Send("下载失败")
	}

	zap.S().Debugf("Finished command with output\n %v", outStr)
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
