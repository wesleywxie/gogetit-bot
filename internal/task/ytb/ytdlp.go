package ytb

import (
	"fmt"
	"github.com/wesleywxie/gogetit/internal/config"
	"github.com/wesleywxie/gogetit/internal/task"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v3"
	"os/exec"
	"strings"
)

func GetFilename(c tb.Context, msg *tb.Message, url string, gen chan string) {
	var filename string

	cmd := exec.Command("yt-dlp",
		"--print", "filename",
		"--output", "%(title)s.%(ext)s",
		"--trim-filenames", "50",
		url,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		zap.S().Warnw("Failed to extract filename",
			"url", url,
			"error", err.Error(),
		)
		_, _ = c.Bot().Edit(msg, "下载失败")
	} else {
		filename = strings.TrimSuffix(string(out), "\n")
	}
	gen <- filename
}


func ExecDownload(c tb.Context, msg *tb.Message, url string, gen chan string, download chan string) {

	filename := <-gen

	if len(filename) > 0 {
		_, _ = c.Bot().Edit(msg, fmt.Sprintf("正在下载 [%v]", filename))

		args := make([]string, 0, 7)
		args = append(args, "--downloader", "aria2c")
		args = append(args, "--downloader-args", fmt.Sprintf("-x %d -k 1M", config.MaxThreadNum))
		args = append(args, "--output", filename)
		if len(config.OutputDir) > 0 {
			args = append(args, "--paths", config.OutputDir)
		}
		if len(config.CookieFile) > 0 {
			args = append(args, "--cookies", config.CookieFile)
		}
		if len(config.UserAgent) > 0 {
			args = append(args, "--user-agent", config.UserAgent)
		}
		args = append(args, url)

		command := "yt-dlp"

		err := task.Proceed(command, args...)

		if err != nil {
			zap.S().Warnw("Failed to download",
				"url", url,
				"error", err.Error(),
			)
		}

		_, _ = c.Bot().Edit(msg, fmt.Sprintf("下载完成 [%v]", filename))
	}

	download <- filename
}
