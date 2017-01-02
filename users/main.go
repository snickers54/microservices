package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/snickers54/microservices/library/models"
	"github.com/snickers54/microservices/library/utils"
	"github.com/spf13/viper"
)

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage: %s <config-filepath> \nArguments:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  config-filepath: Path to the config file for the API..\n")
}

func createService() models.Service {
	return models.Service{
		Name: viper.GetString("service.name"),
		Port: viper.GetString("service.port"),
		Version: models.Version{
			Name:  viper.GetString("service.version.name"),
			Major: uint64(viper.GetInt64("service.version.major")),
			Minor: uint64(viper.GetInt64("service.version.minor")),
			Patch: uint64(viper.GetInt64("service.version.patch")),
		},
	}
}

func main() {
	if len(os.Args) != 2 {
		log.WithField("args", os.Args).Fatal("Number of arguments invalid.")
		usage()
	}
	utils.InitConfig(os.Args[1])
	service := createService()
	log.WithField("service", service).Info("Creating service object.")
	addr := fmt.Sprintf("%s://%s/services", viper.GetString("gateway.protocol"), viper.GetString("gateway.addr"))
	log.WithField("gateway.dsn", addr).Info("Registering to gateway.")
	if errs := service.Register(addr); len(errs) > 0 {
		log.WithField("errors", errs).Fatal("Couldn't register to the gateway.")
	}

}