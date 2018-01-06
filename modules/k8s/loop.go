package k8s

import (
	"context"

	"time"

	"strings"

	"sort"

	"encoding/json"

	"clickyab.com/cluster-tools/modules/k8s/models"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
)

type nodes []models.Node

func (n nodes) Len() int {
	return len(n)
}

func (n nodes) Less(i, j int) bool {
	return strings.Compare(n[i].Name, n[j].Name) < 0
}

func (n nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

type looper struct {
	domains []string
	nodes   nodes
}

func getSubDomain(s string) (string, string) {
	parts := strings.Split(strings.TrimSpace(s), ".")
	if len(s) < 2 {
		return "", s
	}
	if len(parts) == 2 {
		return "@", strings.Join(parts, ".")
	}

	return strings.Join(parts[:len(parts)-2], "."), strings.Join(parts[len(parts)-2:], ".")
}

func (l *looper) Initialize(ctx context.Context) {
	safe.ContinuesGoRoutine(ctx, func(c context.CancelFunc) { l.loop(ctx, c) }, 10*time.Second)
}

func (l *looper) checkDomain(s []string) bool {
	sort.Strings(s)
	defer func() {
		l.domains = s
	}()
	if len(s) != len(l.domains) {
		return false
	}
	for i := range s {
		if l.domains[i] != s[i] {
			return false
		}
	}

	return true
}

func (l *looper) checkNodes(n nodes) bool {
	sort.Sort(n)
	defer func() {
		l.nodes = n
	}()
	if len(n) != len(l.nodes) {
		return false
	}
	for i := range n {
		if l.nodes[i].Active != n[i].Active {
			return false
		}
	}

	return true
}

func (l *looper) reSync() {
	domains := map[string][]string{}
	for i := range l.domains {
		s, d := getSubDomain(l.domains[i])
		domains[d] = append(domains[d], s)
	}

	var ips []string
	for i := range l.nodes {
		if l.nodes[i].Active {
			ips = append(ips, l.nodes[i].IP)
		}
	}

	m := map[string]interface{}{
		"ips":     ips,
		"domains": domains,
	}
	d, err := json.MarshalIndent(m, "", "\t")
	assert.Nil(err)
	logrus.Debugf(string(d))
	models.RefreshDNS(ips, domains)
}

func (l *looper) loop(ctx context.Context, cnl context.CancelFunc) {
	d := l.checkDomain(models.Domains())
	n := l.checkNodes(nodes(models.GetNodes()))

	if !d || !n {
		l.reSync()
	}
}

func init() {
	initializer.Register(&looper{}, 0)
}
