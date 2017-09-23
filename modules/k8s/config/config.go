package kcfg

import "github.com/clickyab/services/config"

// BlackList list of nodes we dont need to check comma separated
var BlackList = config.RegisterString("nodemon.modules.kuber.black", "kube-0.clickyab.ae", "do not check these nodes")
