package util

import (
	"fmt"
	"github.com/wesleywxie/gogetit/internal/config"
)

func init() {
	clientInit()
}

func BuildYtdlpArgs(url, filename string) []string {
	args := make([]string, 0, 7)
	args = append(args, "--downloader", "aria2c")
	args = append(args, "--downloader-args", fmt.Sprintf("-x %d -k 1M", config.MaxThreadNum))
	args = append(args, "-o", filename)
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

func GetYtdlpFilename(url string) []string {
	args := []string {
		"--print", "filename",
		"--o", "%(title)s.%(ext)s",
		"--restrict-filenames",
		"--trim-filenames", "50",
		url,
	}
	return args
}