package server

// use stddev to spot it becoming unhealthy

import (
	"time"
)

type heartbeat struct {
	last    time.Time
	maxDiff time.Duration
}

func newHeartbeat(maxDiff time.Duration) *heartbeat {
	return &heartbeat{
		last:    time.Now(),
		maxDiff: maxDiff,
	}
}

func (self *heartbeat) beat() {
	self.last = time.Now()
}

func (self *heartbeat) healthy() bool {
	if t := self.last.Add(self.maxDiff); t.After(time.Now()) {
		return true
	}

	return false
}
