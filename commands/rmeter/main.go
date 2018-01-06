package main

import (
	_ "clickyab.com/cluster-tools/modules/metrics"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/shell"
	"github.com/sirupsen/logrus"

	"gopkg.in/fzerorubigd/onion.v3"
)

func main() {
	l := onion.NewDefaultLayer()
	config.Initialize("clickyab", "rmeter", "RMR", l)
	defer initializer.Initialize()()

	sig := shell.WaitExitSignal()
	logrus.WithField("signal", sig).Debug("Exit")
}
