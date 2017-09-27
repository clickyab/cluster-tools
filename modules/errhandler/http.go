package errhandler

import (
	"github.com/clickyab/services/framework/router"
)

type initRouter struct {
}

// Ignoring the mount path
func (initRouter) Routes(mux router.Mux) {
	mux.GET("/healthz", healthz)
	mux.GET("/", errCheck)
}

func init() {
	router.Register(&initRouter{})
}
