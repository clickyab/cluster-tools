package kcfg

import "github.com/clickyab/services/config"

var BlackList = config.RegisterString("nodemon.modules.kuber.black", "kube-0.clickyab.ae,kube-20.clickyab.ae", "do not check these nodes")
