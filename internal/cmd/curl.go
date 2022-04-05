package cmd

import "github.com/wesleywxie/gogetit-bot/internal/config"

func CheckLiveness(url string) (bool, error) {
	output, error := proceed("curl",
		"-A", config.UserAgent,
		"| grep -o", "https://edge[0-9]*.stream.highwebmedia.com.*/playlist.m3u8",
		"-s", url)

	return len(output) > 0, error
}
