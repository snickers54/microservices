package main

import (
	"fmt"
	"os"

	"github.com/snickers54/microservices/gateway/handlers"
	"github.com/snickers54/microservices/gateway/network"
	"github.com/spf13/viper"
)

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage: %s <config-filepath> \nArguments:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  config-filepath: Path to the config file for the API..\n")
	os.Exit(0)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	InitConfig(os.Args[1])
	network.InitCluster(viper.GetString("cluster.name"))
	network.InsertEndpoint("test", &network.Service{})
	handlers.Start()
	version := network.Version{}
	version.Parse("-wfsdghs best1.3.2ever dsgfhs-")
}
