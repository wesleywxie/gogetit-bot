package model

import (
	"github.com/jinzhu/gorm"
	"github.com/wesleywxie/gogetit/internal/log"
	"moul.io/zapgorm"
	"github.com/wesleywxie/gogetit/internal/config"
	"go.uber.org/zap"
)

var db *gorm.DB

// InitDB init db object
func InitDB() {
	connectDB()
	configDB()
	updateTable()
}

func connectDB() {
	if config.RunMode == config.TestMode {
		return
	}

	var err error
	db, err = gorm.Open("sqlite3", config.SQLitePath)
	if err != nil {
		zap.S().Fatalf("connect db failed, err: %+v", err)
	}
}

func configDB() {
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(50)
	db.LogMode(config.DBLogMode)
	db.SetLogger(zapgorm.New(log.Logger.WithOptions(zap.AddCallerSkip(7))))
}

func updateTable() {
	// Put tables here
}
