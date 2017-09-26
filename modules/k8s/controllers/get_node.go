package controllers

import (
	"strings"

	"clickyab.com/cluster-tools/modules/k8s/config"
	"github.com/clickyab/services/array"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
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

var k8sClient *kubernetes.Clientset

func init() {
	config, err := rest.InClusterConfig()
	assert.Nil(err)
	k8sClient, err = kubernetes.NewForConfig(config)
	assert.Nil(err)

}

// GetNodes return nodes array with their statuses
func GetNodes() []Node {
	var currentNode []Node

	nodes, err := k8sClient.CoreV1().Nodes().List(metav1.ListOptions{})
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
		var node = Node{
			Name:   n.Name,
			IP:     internalIP,
			Status: false,
		}
		for _, s := range n.Status.Conditions {
			if s.Type == v1.NodeConditionType("Ready") && s.Status == v1.ConditionStatus("True") {
				node.Status = true
				break
			}
		}
		currentNode = append(currentNode, node)

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

var domainBlackList = config.RegisterString("kub.domains.blacklist", "", "comma separate domains name")

// Domains return all domains from ingress
func Domains() []string {
	bl := strings.Split(domainBlackList.String(), ",")
	h := make(map[string]int)
	ns, err := k8sClient.CoreV1().Namespaces().List(metav1.ListOptions{})
	assert.Nil(err)
	for _, n := range ns.Items {
		ng := k8sClient.ExtensionsV1beta1().Ingresses(n.Name)

		a, err := ng.List(metav1.ListOptions{})
		assert.Nil(err)
		for _, q := range a.Items {
			for _, ru := range q.Spec.Rules {
				if array.StringInArray(ru.Host, bl...) {
					continue
				}
				h[ru.Host] = 1
			}
		}
	}
	res := make([]string, 0)
	for k := range h {
		res = append(res, k)
	}
	return res
}
