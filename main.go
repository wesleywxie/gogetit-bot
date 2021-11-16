package main

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/wesleywxie/gogetit/internal/log"
	"github.com/wesleywxie/gogetit/internal/model"
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
	go handleSignal()
}

func handleSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-c

	gracefullyShutdown()
	os.Exit(0)
}


func gracefullyShutdown() {
	zap.S().Infof("Shutting down gracefully...")
}
