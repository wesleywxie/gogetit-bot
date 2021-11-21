package task

import (
	"fmt"
	"github.com/wesleywxie/gogetit/internal/config"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v3"
	"os"
	"path/filepath"
)

func Sync(c tb.Context, msg *tb.Message, download chan string) {
	filename := <- download
	if config.AutoUpload {
		_, _ = c.Bot().Edit(msg, fmt.Sprintf("正在上传 [%v]", filename))
		file := filepath.Join(config.OutputDir, filename)
		command := "rclone"
		args := []string{
			"move", "--ignore-existing",
			file,
			fmt.Sprintf("%s:upload/2021-11-21", config.AutoUploadDrive),
		}

		err := Proceed(command, args...)

		if err != nil {
			_, _ = c.Bot().Edit(msg, fmt.Sprintf("上传失败 [%v]", filename))
			zap.S().Warnw("Failed to sync",
				"filename", filename,
				"error", err.Error(),
			)
			_ = os.Remove(file)
		}
		_, _ = c.Bot().Edit(msg, fmt.Sprintf("上传成功 [%v]", filename))
	}
}