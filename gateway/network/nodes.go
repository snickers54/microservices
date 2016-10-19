package network

import "fmt"

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
	for _, currentNode := range *self {
		if currentNode.IP == node.IP && currentNode.Port == node.Port {
			return false
		}
	}
	*self = append(*self, node)
	return true
}
