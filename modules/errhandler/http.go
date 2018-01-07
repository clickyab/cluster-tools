package errhandler

import (
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
)

type initRouter struct {
}

// Ignoring the mount path
func (initRouter) Routes(mux framework.Mux) {
	mux.GET("root", "/", errCheck)
}

func init() {
	router.Register(&initRouter{})
}
