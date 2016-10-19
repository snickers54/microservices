package main

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// parse config-filepath to be decomposed on something we can pass to viper
func explodeConfigPath(path string) (dir, name, ext string) {
	dir, name = filepath.Split(path)
	ext = filepath.Ext(path)
	name = strings.TrimRight(name, ext)
	return dir, name, ext
}

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
}

func InitConfig(path string) error {
	dir, name, _ := explodeConfigPath(path)
	viper.SetConfigName(name)
	if len(dir) > 0 {
		viper.AddConfigPath(dir)
	}
	err := viper.ReadInConfig()
	return err
}
