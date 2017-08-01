package slavemonitoring

import (
	"context"
	"net/http"
	"regexp"
	"strconv"

	"github.com/clickyab/services/framework/router"
	"github.com/clickyab/services/mysql"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
)

var (
	regex = regexp.MustCompile(`db(\d+)`)
)

type route struct{}

func (route) Routes(mux *xmux.Mux, moountPoint string) {
	mux.GET("/status/:dbnum", xhandler.HandlerFuncC(monitor))
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
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func init() {
	router.Register(route{})
}
