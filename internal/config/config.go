package config

import (
	"fmt"
	"github.com/spf13/viper"
)


type RunType string

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	ProjectName          string = "gogetit"
)

func AppVersionInfo() (s string) {
	s = fmt.Sprintf("version %v, commit %v, built at %v", version, commit, date)
	return
}

// GetString get string config value by key
func GetString(key string) string {
	var value string
	if viper.IsSet(key) {
		value = viper.GetString(key)
	}

	return value
}
