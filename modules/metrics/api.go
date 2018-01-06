package metrics

import (
	"context"
	"encoding/json"
	"net/http"

	"fmt"

	"net/url"

	"time"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/initializer"
)

var (
	httpClient *http.Client
	endpoint   *url.URL
)

type rmqManagement struct {
	Name                      string `json:"name"`
	Vhost                     string `json:"vhost"`
	Durable                   bool   `json:"durable"`
	AutoDelete                bool   `json:"auto_delete"`
	Exclusive                 bool   `json:"exclusive"`
	Node                      string `json:"node"`
	DiskWrites                int    `json:"disk_writes"`
	DiskReads                 int    `json:"disk_reads"`
	MessagesPersistent        int    `json:"messages_persistent"`
	MessagesUnacknowledgedRAM int    `json:"messages_unacknowledged_ram"`
	MessagesReadyRAM          int    `json:"messages_ready_ram"`
	MessagesRAM               int    `json:"messages_ram"`
	State                     string `json:"state"`
	Memory                    int    `json:"memory"`
	Consumers                 int    `json:"consumers"`
	MessagesUnacknowledged    int    `json:"messages_unacknowledged"`
	MessagesReady             int    `json:"messages_ready"`
	Messages                  int    `json:"messages"`
	Reductions                int    `json:"reductions"`
}

func getStatus(vHost, queue string) (*rmqManagement, error) {
	if vHost == "/" || vHost == "" {
		vHost = "%2f"
	}
	path := endpoint.String() + "/" + vHost + "/" + queue
	r, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	nCtx, _ := context.WithTimeout(context.Background(), time.Second)
	resp, err := httpClient.Do(r.WithContext(nCtx))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("invalid status %d", resp.StatusCode)
		return nil, err
	}

	mgm := rmqManagement{}
	err = json.NewDecoder(resp.Body).Decode(&mgm)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return &mgm, nil
}

type initClient struct {
}

func (*initClient) Initialize(context.Context) {
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
			MaxIdleConns:        11,
		},
	}
	var err error
	endpoint, err = url.Parse(dsn.String())
	assert.Nil(err)
	endpoint.Path = "/api/queues"
}

func init() {
	initializer.Register(&initClient{}, 100)
}
