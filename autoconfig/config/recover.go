package config

import "thingworks/common/utils"

type RecoverConfig struct {
	BidderExpiryTime int
	Cron             string
	Limit            int
}

func (recoverConfig RecoverConfig) Check() error {
	return utils.CheckCron(recoverConfig.Cron)
}

