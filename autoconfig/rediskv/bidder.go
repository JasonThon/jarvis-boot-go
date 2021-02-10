package rediskv

import (
	"log"
)

type Bidder interface {
	TryBid(string, int) BidResult
	TryBidAndRun(string, int, func())
}

type BidResult struct {
	Success bool
}

type DistributedRightBidder struct {
	conn RedisConnection
}

func (bidder *DistributedRightBidder) TryBidAndRun(key string, expireSeconds int, callback func()) {
	result := bidder.TryBid(key, expireSeconds)

	if result.Success {
		callback()
	} else {
		log.Printf("Failed to bid for key [%s]", key)
	}
}

func (bidder *DistributedRightBidder) TryBid(key string, expireSeconds int) BidResult {
	err := bidder.conn.SetIfNotExistWithExpiryTime(key, 1, expireSeconds)

	if err != nil {
		return BidResult{
			Success: false,
		}
	}

	return BidResult{
		Success: true,
	}
}
