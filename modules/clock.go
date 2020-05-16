package modules

import (
	"time"
)

type ClockModule struct {
	// Mon Jan 2 15:04:05 -0700 MST 2006
	layout string
	Name   string
}

func Clock() *ClockModule {
	clock := ClockModule{"Mon 02, Jan 15:04", "local_time"}
	return &clock
}

func (c *ClockModule) Output() string {
	now := time.Now()
	return SimpleOutput(now.Format(c.layout))
}
