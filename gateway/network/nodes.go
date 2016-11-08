package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	log.Println("Update node from remoteAddr : ", remoteAddr)
	var err error
	self.IP, self.Port, err = context.ExtractIpPort(remoteAddr)
	if err != nil {
		// maybe not panicking ?
		panic("RemoteAddr doesn't seems valid ?! See by yourself : " + remoteAddr)
	}
}

func (self *Node) registerToMaster(ip, port string) {
	route := fmt.Sprintf("http://%s:%s%s", ip, port, "/cluster/nodes")
	jsonBytes, _ := json.Marshal(*self)
	req, err := http.NewRequest("POST", route, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Duration(5 * time.Second)}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	node := Node{}
	err = json.NewDecoder(resp.Body).Decode(&node)
	resp.Body.Close()
	if err != nil {
		log.Println("Register to Master :", err)
	}
	GetCluster().Nodes.Me().IP = node.IP
	GetCluster().Nodes.Add(Node{
		IP:     ip,
		Port:   port,
		Status: STATUS_ACTIVE,
		Myself: false,
	})
	// io.Copy(ioutil.Discard, resp.Body)
}

func (self *Node) ping() {
	route := fmt.Sprintf("http://%s:%s%s", self.IP, self.Port, "/healthcheck")
	req, err := http.NewRequest("GET", route, nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Duration(5 * time.Second)}
	resp, err := client.Do(req)
	self.Status = STATUS_ACTIVE
	if err != nil {
		self.Status = STATUS_INACTIVE
		return
	}
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
}

func (self *Node) Start() {
	log.Println("Node routines here !")
	// if slave config then contact master to register as cluster node
	if viper.IsSet("cluster.master") {
		log.Println("Apparently I'm a slave, I've to register to the known master")
		self.registerToMaster(viper.GetString("cluster.master.ip"), viper.GetString("cluster.master.port"))
	}
	log.Println("Entering infinite loop to deal with pinging other nodes !")
	for {
		for _, node := range GetCluster().Nodes {
			if node.Myself == false {
				log.Println("Trying to ping : ", node)
				go node.ping()
			}
		}
		log.Println("Now waiting for ", viper.GetDuration("cluster.ping.interval"), " seconds")
		time.Sleep(viper.GetDuration("cluster.ping.interval"))
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

func (self *Nodes) Me() *Node {
	for _, currentNode := range *self {
		if currentNode.Myself == true {
			return &currentNode
		}
	}
	return nil
}
