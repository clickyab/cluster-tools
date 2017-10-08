package main

import (
	_ "clickyab.com/cluster-tools/modules/errhandler"
	"github.com/clickyab/services/config"
	_ "github.com/clickyab/services/healthz"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/shell"
	"github.com/sirupsen/logrus"
)

const (
	org    string = "clickyab"
	app           = "default-http-backend"
	prefix        = "DHB"
)

func main() {
	config.Initialize(org, app, prefix)
	defer initializer.Initialize()()

	sig := shell.WaitExitSignal()
	logrus.Infof("Signal %s received, Exiting ...", sig)
}
