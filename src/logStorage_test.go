package raftlib

import (
	"encoding/json"
	"log"
	"testing"
)

func TestLogStorageInitialLogStorage(t *testing.T) {
	testStorage()
}

// 测试结构体byte化储存
func testStorage() {
	logStorage := LogStorage{}
	logEntry := Log{
		Index:   1,
		Term:    1,
		LogType: LogCommand,
		Data:    []byte("set key val"),
	}
	logEntryEncoded, err := json.Marshal(logEntry)
	if err != nil {
		log.Fatal("创建logEntryEncoded失败.\n")
	}
	logStorage.entries = append(logStorage.entries, logEntryEncoded)
	logEntryInLogStorage := logStorage.entries[0]
	var logEntryInLogStorageDecoded Log
	err = json.Unmarshal(logEntryInLogStorage, &logEntryInLogStorageDecoded)
	if err != nil {
		log.Fatal("创建logEntryInLogStorageDecoded失败.\n")
	}
	log.Println(logEntryInLogStorageDecoded)
	decodedData := []byte(logEntryInLogStorageDecoded.Data)
	log.Println(string(decodedData))
}
