package config

import "thingworks.net/thingworks/jarvis-boot/utils"

type RecoverConfig struct {
	BidderExpiryTime int
	Cron             string
	Limit            int
}

func (recoverConfig RecoverConfig) Check() error {
	return utils.CheckCron(recoverConfig.Cron)
}
