package example

import (
	raftlib "raftlib/src"
	"time"
)

type RaftConfig struct {
	Config      *raftlib.Config
	LocalServer raftlib.Server
	// redis连接设置
	REDIS_ADDRESS      string
	REDIS_PASSWORD     string
	REDIS_DB           int
	REDIS_DIAL_TIMEOUT time.Duration
	REDIS_READ_TIMEOUT time.Duration
}

func CreateConfig() *RaftConfig {
	config := &RaftConfig{
		Config: raftlib.CreateConfig(),
		// redis连接设置
		REDIS_ADDRESS:      "127.0.0.1:6379",
		REDIS_PASSWORD:     "",
		REDIS_DB:           0,
		REDIS_DIAL_TIMEOUT: time.Second * 2,
		REDIS_READ_TIMEOUT: time.Second * 2,
	}
	config.Config.DEBUG = false
	// 配置Servers
	config.Config.Servers = []raftlib.Server{
		{raftlib.Voter, "0", "192.168.224.101", "5678"},
		{raftlib.Voter, "1", "192.168.224.102", "5678"},
		{raftlib.Voter, "2", "192.168.224.103", "5678"},
	}
	// 创建client
	config.Config.LocalServer = config.Config.Servers[2]
	config.LocalServer = config.Config.Servers[2]
	return config
}
