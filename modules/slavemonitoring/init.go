package slavemonitoring

import (
	"context"
	"net/http"
	"regexp"
	"strconv"

	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
	"github.com/clickyab/services/mysql"
	"github.com/rs/xmux"
)

var (
	regex = regexp.MustCompile(`db(\d+)`)
)

type route struct{}

func (r route) Routes(mux framework.Mux) {
	mux.GET("status", "/status/:dbnum", r.monitor)
}

func (*route) monitor(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)

}

func init() {
	router.Register(route{})
}
