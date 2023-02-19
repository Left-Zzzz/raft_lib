package raftlib

import (
	"math/rand"
	"time"
)

// randomTimeout : 返回介于minVal和两倍minVal的结果.
func randomTimeout(minVal time.Duration) <-chan time.Time {
	if minVal == 0 {
		return nil
	}
	extra := time.Duration(rand.Int63()) % minVal
	return time.After(minVal + extra)
}
