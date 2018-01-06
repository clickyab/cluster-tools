package metrics

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

var (
	dsn     = config.RegisterString("metrics.amqp.management", "http://cluster-tools:bita123@127.0.0.1:15672", "user password are required")
	metrics = config.RegisterString("metrics.queue_names", "", "comma separated list of vhost/queue")
)

type route struct {
	queue []struct {
		VHost, Queue string
	}
}

func (rr *route) Routes(mux framework.Mux) {
	mux.GET("metrics", "/metrics", rr.monitor)
	all := strings.Split(metrics.String(), ",")
	for i := range all {
		v := strings.Split(all[i], "/")
		if len(v) != 2 {
			continue
		}
		rr.queue = append(rr.queue, struct{ VHost, Queue string }{VHost: strings.Trim(v[0], "\n\t\r "), Queue: strings.Trim(v[1], "\n\t\r ")})
	}
}

func (rr *route) monitor(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	for i := range rr.queue {
		data, err := getStatus(rr.queue[i].VHost, rr.queue[i].Queue)
		if err != nil {
			fmt.Fprintf(buf, "# %s/%s Err!: %s\n", rr.queue[i].VHost, rr.queue[i].Queue, err)
			continue
		}
		fmt.Fprintf(buf, "# %s/%s Consumers : %d Messages : %d\n", rr.queue[i].VHost, rr.queue[i].Queue, data.Consumers, data.Messages)
		fmt.Fprintf(buf, "%s_%s %d\n", rr.queue[i].VHost, rr.queue[i].Queue, data.Messages)
	}
	fmt.Fprintf(w, buf.String())
}

func init() {
	router.Register(&route{})
}
