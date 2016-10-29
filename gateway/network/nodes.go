package network

import (
	"fmt"

	"github.com/snickers54/microservices/gateway/context"
)

type Node struct {
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Port   string `json:"port"`
	Status string `json:"status"`
}

type Nodes []Node

func (self *Node) String() string {
	return fmt.Sprintf("%s(%s) - %s:%s", self.Name, self.Status, self.IP, self.Port)
}

func (self *Nodes) Add(node Node) bool {
	if self.Exists(node) == true {
		return false
	}
	*self = append(*self, node)
	return true
}

func (self *Nodes) Exists(node Node) bool {
	for _, currentNode := range *self {
		if currentNode.IP == node.IP && currentNode.Port == node.Port {
			return true
		}
	}
	return false
}

func (self *Node) UpdateFromAddr(remoteAddr string) {
	var err error
	self.IP, self.Port, err = context.ExtractIpPort(remoteAddr)
	if err != nil {
		// maybe not panicking ?
		panic("RemoteAddr doesn't seems valid ?! See by yourself : " + remoteAddr)
	}
}
