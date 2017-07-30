package slavemonitoring

import (
	"context"
	"net/http"
	"regexp"

	"strings"

	"github.com/clickyab/services/mysql"

	"strconv"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework/router"
	"github.com/clickyab/services/initializer"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
)

var (
	regex       = regexp.MustCompile(`db(/d+)`)
	rdsnSlice   = config.RegisterString("services.mysql.rdsn", "root:bita123@tcp(127.0.0.1:3306)/?charset=utf8&parseTime=true", "comma separated read database dsn")
	connections []string
)

type route struct{}

func (route) Initialize(ctx context.Context) {
	connections = strings.Split(rdsnSlice.String(), ",")
}

func (route) Routes(mux *xmux.Mux, moountPoint string) {
	mux.GET("/healthz", xhandler.HandlerFuncC(healthz))
	mux.GET("/:dbnum/healthz", xhandler.HandlerFuncC(monitor))
}

func healthz(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func monitor(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	dbNum := xmux.Param(ctx, "dbnum")
	slice := regex.FindStringSubmatch(dbNum)
	if len(slice) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbIndex, err := strconv.Atoi(slice[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = mysql.PingDB(dbIndex)
	if err != nil {
		w.Write([]byte(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func init() {
	initializer.Register(route{}, 1)
	router.Register(route{})
}
