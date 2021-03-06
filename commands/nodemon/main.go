package main

import (
	_ "clickyab.com/cluster-tools/modules/k8s"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/shell"
	"github.com/sirupsen/logrus"
)

const (
	org    string = "clickyab"
	app           = "nodemon"
	prefix        = "NM"
)

func main() {
	config.Initialize(org, app, prefix)
	defer initializer.Initialize()()

	sig := shell.WaitExitSignal()
	logrus.Infof("Signal %s received, Exiting ...", sig)
}
