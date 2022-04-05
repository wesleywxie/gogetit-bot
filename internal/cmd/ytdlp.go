package cmd

import (
	"fmt"
	"github.com/wesleywxie/gogetit-bot/internal/config"
	"github.com/wesleywxie/gogetit-bot/internal/model"
	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"
	"os/exec"
	"strings"
	"time"
)

func GetFilename(c tb.Context, msg *tb.Message, url string, gen chan string) {
	var filename string

	cmd := exec.Command("yt-dlp",
		"--print", "filename",
		"--output", "%(title)s.%(ext)s",
		"--trim-filenames", "50",
		"--no-warnings",
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
		outputs := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
		filename = outputs[len(outputs)-1]
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

		_, err := proceed(command, args...)

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

func Recording(subscription model.Subscription, record chan string) {

	now := time.Now()
	filename := fmt.Sprintf("%v.mp4", now.Format("20060201150405"))

	zap.S().Infow("Recording...",
		"url", subscription.Link,
		"filename", filename)
	_, _ = model.UpdateStreamingStatus(subscription.ID, true)

	args := make([]string, 0, 7)
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
	args = append(args, subscription.Link)

	command := "yt-dlp"

	_, err := proceed(command, args...)

	if err != nil {
		zap.S().Warnw("Failed to download",
			"url", subscription.Link,
			"error", err.Error(),
		)
	}
	_, _ = model.UpdateStreamingStatus(subscription.ID, false)

	record <- filename
}
