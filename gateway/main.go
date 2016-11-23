package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/snickers54/microservices/gateway/handlers"
	"github.com/snickers54/microservices/gateway/network"
)

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage: %s <config-filepath> \nArguments:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  config-filepath: Path to the config file for the API..\n")
}

func main() {
	if len(os.Args) != 2 {
		log.WithField("args", os.Args).Fatal("Number of arguments invalid.")
		usage()
	}
	InitConfig(os.Args[1])
	network.InitCluster()
	handlers.Start()
}
