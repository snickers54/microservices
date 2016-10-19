package network

import (
	"fmt"

	"github.com/spf13/viper"
)

type Cluster struct {
	Name     string    `json:"name"`
	Nodes    Nodes     `json:"nodes"`
	Services Services  `json:"services"`
	Map      radixTree `json:"map"`
}

func (self *Cluster) String() string {
	return fmt.Sprintf("%s contains %d nodes.", self.Name, len(self.Nodes))
}

var _cluster *Cluster

func InitCluster(name string) {
	_cluster = new(Cluster)
	_cluster.Name = name
	_cluster.Map = radixTree{_map}
}

func GetCluster() *Cluster {
	if _cluster == nil {
		InitCluster(viper.GetString("cluster.name"))
	}
	return _cluster
}
