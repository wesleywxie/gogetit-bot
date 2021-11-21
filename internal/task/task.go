package task

import (
	"bytes"
	"go.uber.org/zap"
	"io"
	"os/exec"
)

var (
	taskList []Task
)

type Task interface {
	Start()
	Stop()
	Name() string
}

func init() {
	taskList = []Task{}
}

func registerTask(task Task) {
	taskList = append(taskList, task)
}

func StartTasks() {
	for _, task := range taskList {
		zap.S().Infof("Start task %v", task.Name())
		task.Start()
	}
}

func StopTasks() {
	for _, task := range taskList {
		zap.S().Infof("Stop task %v", task.Name())
		task.Stop()
	}
}

func Proceed(command string, args... string) error {
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
