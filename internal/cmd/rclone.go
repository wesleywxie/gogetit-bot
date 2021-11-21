package cmd

import (
	"fmt"
	"github.com/wesleywxie/gogetit/internal/config"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v3"
	"os"
	"path/filepath"
	"time"
)

func Sync(c tb.Context, msg *tb.Message, download chan string) {
	filename := <- download
	if config.AutoUpload {

		now := time.Now()
		dir := fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())

		_, _ = c.Bot().Edit(msg, fmt.Sprintf("正在上传 [%v]", filename))
		file := filepath.Join(config.OutputDir, filename)
		command := "rclone"
		args := []string {
			"move", "--ignore-existing", file,
			fmt.Sprintf("%s:%s/%s", config.AutoUploadDrive, config.ProjectName, dir),
		}

		err := proceed(command, args...)

		if err != nil {
			_, _ = c.Bot().Edit(msg, fmt.Sprintf("上传失败 [%v]", filename))
			zap.S().Warnw("Failed to sync",
				"filename", filename,
				"error", err.Error(),
			)
			// Delete the downloaded file no matter what
			_ = os.Remove(file)
		}
		_, _ = c.Bot().Edit(msg, fmt.Sprintf("上传成功 [%v]", filename))
	}
}