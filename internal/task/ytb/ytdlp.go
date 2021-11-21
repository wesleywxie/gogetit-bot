package ytb

import (
	"fmt"
	"github.com/wesleywxie/gogetit/internal/config"
	"github.com/wesleywxie/gogetit/internal/task"
	"go.uber.org/zap"
	"os/exec"
	"strings"
)

func GetFilename(url string) (filename string, err error) {
	cmd := exec.Command("yt-dlp",
		"--print", "filename",
		"--output", "%(title)s.%(ext)s",
		"--trim-filenames", "50",
		url,
	)
	out, err := cmd.CombinedOutput()
	if err == nil {
		filename = strings.TrimSuffix(string(out), "\n")
	}
	return
}


func ExecDownload(url, filename string) (err error) {
	args := buildDownloadArgs(url, filename)

	zap.S().Debugf("Executing command yt-dlp %v", args)
	cmd := exec.Command("yt-dlp", args...)

	return task.Proceed(cmd)
}


func buildDownloadArgs(url, filename string) []string {
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
	return append(args, url)
}