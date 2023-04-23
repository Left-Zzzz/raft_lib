package main

import (
	"fmt"
	example "raftlib/example"
	raftlib "raftlib/src"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

var serverConfig *example.RaftConfig

// 服务端
type Server struct {
	raft        *raftlib.Raft
	redisClient *redis.Client
}

// 注册回调函数
func (s *Server) registerExecCommandFunc() {
	cbFunc := func(b []byte) ([]byte, error) {
		command := string(b)
		// 解析Redis命令
		commandSplited := strings.Split(command, " ")
		// fmt.Printf("cbFunc running, command:%v, commandSplited:%v\n", command, commandSplited)
		var returnValue []byte
		// 执行Redis命令
		switch strings.ToLower(commandSplited[0]) {
		case "set":
			fmt.Printf("cbFunc: set\n")
			result, err := s.redisClient.Set(commandSplited[1], commandSplited[2], 0).Result()
			returnValue = []byte(result)
			if err != nil {
				fmt.Printf("Set error:%v \n", err)
				return nil, err
			}
			fmt.Printf("Set result: %v\n", result)
		case "del":
			// fmt.Printf("cbFunc: del\n")
			key := strings.Trim(commandSplited[1], " \n")
			result, err := s.redisClient.Del(key).Result()
			returnValue = []byte(strconv.FormatInt(result, 10))
			if err != nil {
				fmt.Printf("Del error:%v \n", err)
				return nil, err
			}
			fmt.Printf("Del result: %v\n", result)
		default:
			returnValue = []byte("")
			// fmt.Printf("Invalid command:%v \n", commandSplited[0])
		}
		return returnValue, nil
	}
	s.raft.RegisterCallBackFunc(cbFunc)
}

func (s *Server) createRedisClient() {
	// 通过配置文件读取Redis客户端参数
	redisClient := redis.NewClient(&redis.Options{
		Addr:        serverConfig.REDIS_ADDRESS,
		Password:    serverConfig.REDIS_PASSWORD,
		DB:          serverConfig.REDIS_DB,
		DialTimeout: serverConfig.REDIS_DIAL_TIMEOUT,
		ReadTimeout: serverConfig.REDIS_READ_TIMEOUT,
	})
	// 测试连接是否成功
	pong, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("createRedisClient(): pong: %v\n", pong)
	s.redisClient = redisClient
}

func createServer(node raftlib.Server) *Server {
	server := &Server{}
	// 创建raft服务
	server.raft = raftlib.CreateRaft(serverConfig.Config)
	// 创建redis客户端
	server.createRedisClient()
	// 注册执行命令回调函数
	server.registerExecCommandFunc()
	// 启动raft服务
	go server.raft.Run()
	fmt.Printf("create server finished.\n")
	return server
}

func main() {
	serverConfig = example.CreateConfig()
	// serverConfig.Config.DEBUG = false

	// 创建server
	createServer(serverConfig.Config.LocalServer)
	for {
		command := ""
		fmt.Scanln(command)
		if command == "stop" {
			break
		}
	}
}
