package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/shell"
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
