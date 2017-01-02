package network

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/snickers54/microservices/library/models"
	"github.com/spf13/viper"
)

type Cluster struct {
	Nodes    Nodes           `json:"nodes"`
	Services models.Services `json:"services"`
	Map      radixTree       `json:"map"`
	Name     string          `json:"name"`
}

func (self *Cluster) String() string {
	return fmt.Sprintf("%s contains %d nodes.", self.Name, len(self.Nodes))
}

var _cluster *Cluster

func InitCluster() {
	_cluster = new(Cluster)
	_cluster.Map = radixTree{_map}
	_cluster.Services = models.Services{}
	_cluster.Nodes = Nodes{}
	_cluster.Name = "Singlette Gateway Cluster prototype"
	node := createOwnNode()
	go node.Start()
	_cluster.Nodes.Add(node)
	log.WithField("cluster", _cluster).Debug("Init cluster !")
}

func createOwnNode() Node {
	return Node{
		Name:   viper.GetString("node.name"),
		Port:   viper.GetString("node.port"),
		Status: models.STATUS_ACTIVE,
		IP:     "localhost",
		Myself: true,
	}
}

func GetCluster() *Cluster {
	if _cluster == nil {
		InitCluster()
	}
	return _cluster
}
