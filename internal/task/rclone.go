package task

import (
	"fmt"
	"github.com/wesleywxie/gogetit/internal/config"
	"os"
	"path/filepath"
)

func Sync(filename string) error {

	file := filepath.Join(config.OutputDir, filename)
	command := "rclone"
	args := []string{
		"move", "--ignore-existing",
		file,
		fmt.Sprintf("%s:upload/2021-11-21", config.AutoUploadDrive),
	}

	err := Proceed(command, args...)

	if err != nil {
		_ = os.Remove(file)
	}

	return err
}