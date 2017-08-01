package main

import (
	_ "clickyab.com/cluster-tools/modules/slavemonitoring"
	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	_ "github.com/clickyab/services/mysql/connection/mysql"
	"github.com/clickyab/services/shell"
)

const (
	org    string = "clickyab"
	app           = "slave-mon"
	prefix        = "SM"
)

func main() {
	config.Initialize(org, app, prefix)
	defer initializer.Initialize()()

	sig := shell.WaitExitSignal()
	logrus.Infof("Signal %s received, Exiting ...", sig)
}
