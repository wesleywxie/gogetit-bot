package task

import (
	"fmt"
	"github.com/wesleywxie/gogetit/internal/config"
	"os"
	"os/exec"
	"path/filepath"
)

func Sync(filename string) error {

	file := filepath.Join(config.OutputDir, filename)

	cmd := exec.Command("rclone",
		"move", "--ignore-existing",
		file,
		fmt.Sprintf("%s:upload/2021-11-21", config.AutoUploadDrive),
		)

	err := Proceed(cmd)

	if err != nil {
		_ = os.Remove(file)
	}

	return err
}