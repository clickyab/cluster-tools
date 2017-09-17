package main

import (
	_ "clickyab.com/cluster-tools/modules/hls/controllers"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	_ "github.com/clickyab/services/kv/redis"
	"github.com/clickyab/services/shell"
	"github.com/sirupsen/logrus"
	"gopkg.in/fzerorubigd/onion.v3"
)

func main() {
	l := onion.NewDefaultLayer()
	l.SetDefault("services.framework.controller.mount_point", "/hls")
	config.Initialize("clickyab", "mantis", "HLS", l)
	defer initializer.Initialize()()

	sig := shell.WaitExitSignal()
	logrus.WithField("signal", sig).Debug("Exit")
}
