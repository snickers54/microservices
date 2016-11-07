package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/snickers54/microservices/gateway/context"
	"github.com/spf13/viper"
)

// DEFINITION OF NODE
type Node struct {
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Port   string `json:"port"`
	Status string `json:"status"`
	Myself bool   `json:"-"`
}

func (self *Node) String() string {
	return fmt.Sprintf("%s(%s) - %s:%s", self.Name, self.Status, self.IP, self.Port)
}

func (self *Node) UpdateFromAddr(remoteAddr string) {
	var err error
	self.IP, self.Port, err = context.ExtractIpPort(remoteAddr)
	if err != nil {
		// maybe not panicking ?
		panic("RemoteAddr doesn't seems valid ?! See by yourself : " + remoteAddr)
	}
}

func registerToMaster(ip, port string) {
	route := fmt.Sprintf("http://%s:%s%s", ip, port, "/cluster/nodes")
	node := Node{
		Name: viper.GetString("node.name"),
		Port: viper.GetString("node.port"),
	}
	jsonBytes, _ := json.Marshal(node)
	req, err := http.NewRequest("POST", route, bytes.NewBuffer(jsonBytes))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	// we have to read all the body to close the connection .. if not it blocks ..
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
}

func (self *Node) Start() {
	// if slave config then contact master to register as cluster node
	if viper.IsSet("cluster.master") {
		registerToMaster(viper.GetString("cluster.master.ip"), viper.GetString("cluster.master.port"))
	}
	for {
		time.Sleep(10 * time.Second)
	}
}

// DEFINITION OF NODES

type Nodes []Node

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
