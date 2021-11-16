package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func init() {
	if isInTests() {
		// 测试环境
		RunMode = TestMode
		return
	}

	workDirFlag := flag.String("d", "./", "work directory of gogetid")
	configFile := flag.String("c", "", "config file of gogetid")
	printVersionFlag := flag.Bool("v", false, "prints gogetid version")

	testing.Init()
	flag.Parse()

	if *printVersionFlag {
		// print version
		fmt.Printf(AppVersionInfo())
		os.Exit(0)
	}

	workDir := filepath.Clean(*workDirFlag)

	if *configFile != "" {
		viper.SetConfigFile(*configFile)
	} else {
		viper.SetConfigFile(filepath.Join(workDir, "config.yml"))
	}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error on reading config file: %s", err))
	}

	if viper.IsSet("sqlite.path") {
		SQLitePath = viper.GetString("sqlite.path")
	} else {
		SQLitePath = filepath.Join(workDir, "data.db")
	}
	log.Println("DB Path:", SQLitePath)
	// 判断并创建SQLite目录
	dir := path.Dir(SQLitePath)
	_, err = os.Stat(dir)
	if err != nil {
		err := os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			log.Printf("mkdir failed![%v]\n", err)
		}
	}

	if viper.IsSet("log.db_log") {
		DBLogMode = viper.GetBool("log.db_log")
	}
}

func isInTests() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test") {
			return true
		}
	}
	return false
}
