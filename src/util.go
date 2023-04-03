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

// 根据传参的index生成后一位index
func genNextLogIndex(index uint32) uint32 {
	// 因为uint32中负一表示为MAX_LOG_INDEX_NUM,所以不能简单+1
	if index == MAX_LOG_INDEX_NUM {
		index = 0
	} else {
		index++
	}
	return index
}

// 根据传参的index生成前一位index
func genPrevLogIndex(index uint32) uint32 {
	// 因为uint32中负一表示为MAX_LOG_INDEX_NUM,所以不能简单-1
	if index == 0 || index == MAX_LOG_INDEX_NUM {
		index = MAX_LOG_INDEX_NUM
	} else {
		index--
	}
	return index
}
