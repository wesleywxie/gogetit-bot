package cmd

import (
	"bytes"
	"go.uber.org/zap"
	"io"
	"os/exec"
)

func proceed(command string, args... string) error {
	cmd := exec.Command(command, args...)

	zap.S().Debugf("Executing command %s with args %v", command, args)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(&stdoutBuf)
	cmd.Stderr = io.MultiWriter(&stderrBuf)

	err := cmd.Run()
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())

	if err != nil {
		zap.S().Debugf("Finished command with error output\n %v", errStr)
	}

	zap.S().Debugf("Finished command with output\n %v", outStr)

	return err
}