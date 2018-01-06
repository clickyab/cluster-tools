package metrics

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
	"github.com/clickyab/services/safe"
	"github.com/streadway/amqp"
)

var (
	dsn      = config.RegisterString("metric.amqp.dsn", "amqp://cluster-tools:bita123@127.0.0.1:5672/", "no vhost please!")
	tryLimit = config.RegisterDuration("services.amqp.try_limit", time.Minute, "")
	metrics  = config.RegisterString("metrics.queue_names", "", "comma separated list of vhost/queue")
)

type route struct {
	conn  map[string]*amqp.Connection
	ch    map[string]*amqp.Channel
	queue []struct {
		VHost, Queue string
	}
}

func (rr *route) Routes(mux framework.Mux) {
	mux.GET("metrics", "/metrics", rr.monitor)
	all := strings.Split(metrics.String(), ",")
	var vHosts []string
	for i := range all {
		v := strings.Split(all[i], "/")
		if len(v) != 2 {
			continue
		}
		vHosts = append(vHosts, strings.Trim(v[0], "\n\t\r "))
		rr.queue = append(rr.queue, struct{ VHost, Queue string }{VHost: strings.Trim(v[0], "\n\t\r "), Queue: strings.Trim(v[1], "\n\t\r ")})
	}

	assert.True(len(vHosts) > 0, "No queue configured for monitoring")
	for i := range vHosts {
		if _, ok := rr.conn[vHosts[i]]; ok {
			continue
		}
		var (
			con *amqp.Connection
			err error
		)
		safe.Try(func() error {
			var err error
			con, err = amqp.Dial(dsn.String() + vHosts[i])
			return err
		}, tryLimit.Duration())
		rr.conn[vHosts[i]] = con
		rr.ch[vHosts[i]], err = con.Channel()
		assert.Nil(err)
	}
}

func (rr *route) monitor(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	for i := range rr.queue {
		ch := rr.ch[rr.queue[i].VHost]
		data, err := ch.QueueInspect(rr.queue[i].Queue)
		if err != nil {
			fmt.Fprintf(buf, "# %s/%s Err!: %s\n", rr.queue[i].VHost, rr.queue[i].Queue, err)
			// Channel is closed because of err, re open it
			rr.ch[rr.queue[i].VHost], err = rr.conn[rr.queue[i].VHost].Channel()
			assert.Nil(err)
			continue
		}
		fmt.Fprintf(buf, "# %s/%s Consumers : %d Messages : %d\n", rr.queue[i].VHost, rr.queue[i].Queue, data.Consumers, data.Messages)
		fmt.Fprintf(buf, "%s_%s %d\n", rr.queue[i].VHost, rr.queue[i].Queue, data.Messages)
	}
	fmt.Fprintf(w, buf.String())
}

func init() {
	router.Register(&route{
		conn: make(map[string]*amqp.Connection),
		ch:   make(map[string]*amqp.Channel),
	})
}
