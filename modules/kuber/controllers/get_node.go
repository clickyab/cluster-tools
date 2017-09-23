package controllers

import (
	"strings"

	"clickyab.com/cluster-tools/modules/kuber/config"
	"github.com/clickyab/services/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Node with ip and status
type Node struct {
	Name   string
	IP     string
	Status bool
	Label  bool
}

// GetNodes return nodes array with their statuses
func GetNodes() []Node {
	var currentNode []Node
	config, err := rest.InClusterConfig()
	assert.Nil(err)
	clientSet, err := kubernetes.NewForConfig(config)
	assert.Nil(err)
	nodes, err := clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	assert.Nil(err)
	for _, n := range nodes.Items {
		var internalIP string
		//check if node exist in config
		if !checkBlacklist(n) {
			continue
		}
		for _, s := range n.Status.Addresses {
			if s.Type == v1.NodeAddressType("InternalIP") {
				internalIP = s.Address
				break
			}
		}
		var stat bool //ready or not
		for _, s := range n.Status.Conditions {
			if s.Type == v1.NodeConditionType("Ready") && s.Status == v1.ConditionStatus("True") {
				stat = true
			}
		}
		if stat { //ready
			currentNode = append(currentNode,
				Node{
					Name:   n.Name,
					IP:     internalIP,
					Status: true,
				})
		} else { //we have problem (not ready node)
			currentNode = append(currentNode,
				Node{
					Name:   n.Name,
					IP:     internalIP,
					Status: false,
				})
		}
	}
	return currentNode
}

// check if the node not exists in the blacklist
func checkBlacklist(node v1.Node) bool {
	blackArr := strings.Split(kcfg.BlackList.String(), ",")
	for i := range blackArr {
		if node.Name == blackArr[i] {
			return false
		}
	}
	return true
}
