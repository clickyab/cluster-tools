package main

import (
	_ "clickyab.com/cluster-tools/modules/errhandler"
	"github.com/clickyab/services/config"
	_ "github.com/clickyab/services/healthz"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/shell"
	"github.com/sirupsen/logrus"
	onion "gopkg.in/fzerorubigd/onion.v3"
)

const (
	org    string = "clickyab"
	app           = "default-http-backend"
	prefix        = "DHB"
)

func main() {
	l := onion.NewDefaultLayer()
	l.SetDefault("services.framework.controller.mount_point", "/")

	config.Initialize(org, app, prefix, l)
	defer initializer.Initialize()()

	sig := shell.WaitExitSignal()
	logrus.Infof("Signal %s received, Exiting ...", sig)
}
