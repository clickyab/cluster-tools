package errhandler

import (
	"github.com/clickyab/services/framework/router"
	"github.com/rs/xmux"

	"github.com/rs/xhandler"
)

type initRouter struct {
}

// Ignoring the mount path
func (initRouter) Routes(mux *xmux.Mux, _ string) {
	mux.GET("/healthz", xhandler.HandlerFuncC(healthz))
	mux.GET("/", xhandler.HandlerFuncC(errCheck))
}

func init() {
	router.Register(&initRouter{})
}
