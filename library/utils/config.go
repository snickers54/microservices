package utils

import (
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/spf13/viper"
)

// parse config-filepath to be decomposed on something we can pass to viper
func explodeConfigPath(path string) (dir, name, ext string) {
	dir, name = filepath.Split(path)
	ext = filepath.Ext(path)
	name = strings.TrimRight(name, ext)
	log.WithFields(log.Fields{
		"folder":    dir,
		"filename":  name,
		"extension": ext,
	}).Debug("Exploding config path.")
	return dir, name, ext
}

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
}

func InitConfig(path string) error {
	log.SetLevel(log.DebugLevel)
	log.WithField("path", path).Debug("Init config.")
	dir, name, _ := explodeConfigPath(path)
	viper.SetConfigName(name)
	if len(dir) > 0 {
		viper.AddConfigPath(dir)
	}
	err := viper.ReadInConfig()
	return err
}
