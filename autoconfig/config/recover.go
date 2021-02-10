package config

import "github.com/thingworks/common/utils"

type RecoverConfig struct {
	BidderExpiryTime int
	Cron             string
	Limit            int
}

func (recoverConfig RecoverConfig) Check() error {
	return utils.CheckCron(recoverConfig.Cron)
}
