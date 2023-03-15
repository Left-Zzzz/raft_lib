package raftlib

import "log"

// 输出调试日志
func logDebug(format string, v ...any) {
	if DEBUG {
		if v != nil {
			log.Printf("[Debug] "+format, v...)
		} else {
			log.Printf("[Debug] " + format)
		}

	}
}

// 输出信息日志
func logInfo(format string, v ...any) {
	if v != nil {
		log.Printf("[Info] "+format, v...)
	} else {
		log.Printf("[Info] " + format)
	}
}

// 输出警告日志
func logWarn(format string, v ...any) {
	if v != nil {
		log.Printf("[Warn] "+format, v...)
	} else {
		log.Printf("[Warn] " + format)
	}
}

// 输出错误日志
func logError(format string, v ...any) {
	if v != nil {
		log.Printf("[Error] "+format, v...)
	} else {
		log.Printf("[Error] " + format)
	}
}
