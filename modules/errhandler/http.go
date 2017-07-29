package errhandler

import (
	"github.com/clickyab/services/framework/router"
	"github.com/rs/xmux"

	"github.com/rs/xhandler"
)

type initRouter struct {
}

func (initRouter) Routes(mux *xmux.Mux, mountPoint string) {
	mux.GET(mountPoint+"/", xhandler.HandlerFuncC(errCheck))
}

func init() {
	router.Register(&initRouter{})
}
