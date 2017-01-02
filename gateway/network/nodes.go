package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"

	"net/http"
	"time"

	"github.com/snickers54/microservices/gateway/context"
	"github.com/snickers54/microservices/library/models"
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
	contextLogger := log.WithField("remote address", remoteAddr)
	contextLogger.Debug("Update node with remote address.")
	var err error
	self.IP, self.Port, err = context.ExtractIpPort(remoteAddr)
	if err != nil {
		contextLogger.Error(err)
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
		log.WithField("ip:port", ip+":"+port).Fatal(err)
	}
	node := Node{}
	err = json.NewDecoder(resp.Body).Decode(&node)
	resp.Body.Close()
	if err != nil {
		log.WithError(err).Fatal("Invalid json response from master.")
	}
	GetCluster().Nodes.Me().IP = node.IP
	GetCluster().Nodes.Add(Node{
		IP:     ip,
		Port:   port,
		Status: models.STATUS_ACTIVE,
		Myself: false,
	})
}

func (self *Node) ping() {
	route := fmt.Sprintf("http://%s:%s%s", self.IP, self.Port, "/healthcheck")
	req, err := http.NewRequest("GET", route, nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Duration(5 * time.Second)}
	resp, err := client.Do(req)
	self.Status = models.STATUS_ACTIVE
	if err != nil {
		log.WithField("node", self).Warn("Node seems unreachable, tagging it as inactive.")
		self.Status = models.STATUS_INACTIVE
		return
	}
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
}

func (self *Node) Start() {
	if viper.IsSet("cluster.master") {
		log.WithField("master", viper.GetStringMap("cluster.master")).Debug("I'm a slave, initiating register call to master.")
		self.registerToMaster(viper.GetString("cluster.master.ip"), viper.GetString("cluster.master.port"))
	}
	log.Debug("Entering async infinite loop of pings")
	for {
		for _, node := range GetCluster().Nodes {
			if node.Myself == false {
				log.WithField("node", node).Debug("Trying to ping.")
				go node.ping()
			}
		}
		time.Sleep(viper.GetDuration("cluster.ping.interval"))
	}
}

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
