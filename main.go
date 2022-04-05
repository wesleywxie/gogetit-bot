package main

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/wesleywxie/gogetit-bot/internal/bot"
	_ "github.com/wesleywxie/gogetit-bot/internal/log"
	"github.com/wesleywxie/gogetit-bot/internal/model"
	"github.com/wesleywxie/gogetit-bot/internal/task"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	// Initialization Script here
	zap.S().Infof("Initialization script...")
}

func main() {
	model.InitDB()
	task.StartTasks()
	go handleShutdownSignal()
	bot.Start()
}

func handleShutdownSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-c

	gracefullyShutdown()
	os.Exit(0)
}

func gracefullyShutdown() {
	task.StopTasks()
	model.Disconnect()
	zap.S().Infof("Shutting down gracefully...")
}
