package kcfg

import "github.com/clickyab/services/config"

// BlackList list of nodes we dont need to check comma separated
var (
	BlackList = config.RegisterString("nodemon.modules.kuber.black", "kube-0.clickyab.ae", "do not check these nodes")
	CFKey     = config.GetString("cluster.cloudflare.key")
	CFMail    = config.GetString("cluster.cloudflare.mail")

	LeastIPNum = config.GetIntDefault("cluster.cloudflare.max_ip_per_dns", 3)
)
