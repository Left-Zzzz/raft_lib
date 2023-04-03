package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	raftlib "raft_lib/src"
	"raftlib/pb"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var config *raftlib.Config = raftlib.CreateConfig()

// 记录本地服务器
var localServer raftlib.Server = config.Servers[0]

// client中redis连接设置
const (
	REDIS_ADDRESS      = "192.168.224.101:6379"
	REDIS_PASSWORD     = ""
	REDIS_DB           = 0
	REDIS_DIAL_TIMEOUT = time.Second * 2
	REDIS_READ_TIMEOUT = time.Second * 2
)

// 客户端
type Client struct {
	clientService pb.RpcServiceClient
	clientConn    *grpc.ClientConn
	leaderService pb.RpcServiceClient
	leaderConn    *grpc.ClientConn
	raft          *raftlib.Raft
	redisClient   *redis.Client
	leader        raftlib.Server
}

// 注册回调函数
func (c *Client) registerExecCommandFunc() {
	cbFunc := func(b []byte) error {
		command := string(b)
		// 解析Redis命令
		commandSplited := strings.Split(command, " ")
		fmt.Printf("cbFunc running, command:%v, commandSplited:%v\n", command, commandSplited)

		// 执行Redis命令
		switch strings.ToLower(commandSplited[0]) {
		case "get":
			fmt.Printf("cbFunc: get\n")
			result, err := c.redisClient.Get(commandSplited[1]).Result()
			if err != nil {
				fmt.Printf("Get error: %v\n", err)
				return err
			}
			fmt.Printf("cbFunc: Get result: %v\n", result)
		case "set":
			fmt.Printf("cbFunc: set\n")
			result, err := c.redisClient.Set(commandSplited[1], commandSplited[2], 0).Result()
			if err != nil {
				fmt.Printf("Set error:%v \n", err)
				return err
			}
			fmt.Printf("Set result: %v\n", result)
		case "del":
			fmt.Printf("cbFunc: del\n")
			result, err := c.redisClient.Del(commandSplited[1], commandSplited[2]).Result()
			if err != nil {
				fmt.Printf("Del error:%v \n", err)
				return err
			}
			fmt.Printf("Del result: %v\n", result)
		default:
			fmt.Printf("Invalid command:%v \n", commandSplited[0])
		}
		return nil
	}
	c.raft.RegisterCallBackFunc(cbFunc)
}

func (c *Client) createRedisClient() {
	// 通过配置文件读取Redis客户端参数
	client := redis.NewClient(&redis.Options{
		Addr:        REDIS_ADDRESS,
		Password:    REDIS_PASSWORD,
		DB:          REDIS_DB,
		DialTimeout: REDIS_DIAL_TIMEOUT,
		ReadTimeout: REDIS_READ_TIMEOUT,
	})
	// 测试连接是否成功
	pong, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("createRedisClient(): pong: %v\n", pong)
	c.redisClient = client
}

func createClient(server raftlib.Server) *Client {
	client := &Client{
		leader: server,
	}
	// 创建raft服务
	config.Localserver = localServer
	client.raft = raftlib.CreateRaft(config)
	// 创建redis客户端
	client.createRedisClient()
	// 注册执行命令回调函数
	client.registerExecCommandFunc()
	// 启动raft服务
	go client.raft.Run()
	// 创建RPC连接
	client.runRpcClient(true)
	client.runRpcClient(false)
	fmt.Printf("createClient finished.\n")
	return client
}

// 启动rpc客户端
func (c *Client) runRpcClient(isLeader bool) error {
	fmt.Printf("runRpcClient, leaderAddress: %v, leaderPort: %v\n", c.leader.Address, c.leader.Port)
	// 创建客户端
	clientConn, err := grpc.Dial(string(c.leader.Address)+":"+string(c.leader.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 断线重连参数
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             100 * time.Millisecond,
			PermitWithoutStream: true}))
	if err != nil {
		log.Fatal(err)
		return err
	}
	// 创建RpcService
	clientService := pb.NewRpcServiceClient(clientConn)
	if isLeader {
		c.leaderConn = clientConn
		c.leaderService = clientService
	} else {
		c.clientConn = clientConn
		c.clientService = clientService
	}

	return nil
}

func (c *Client) execCommand(logType raftlib.LogType, command []byte) (bool, error) {
	fmt.Printf("Client.execCommand.\n")
	ver := raftlib.PROTO_VER_EXEC_COMMAND_REQUEST
	req := &pb.ExecCommandRequest{
		Ver:     &ver,
		LogType: uint32(logType),
		Command: command,
	}
	ctx := context.Background()
	resp, _ := c.leaderService.ExecCommand(ctx, req)
	fmt.Printf("Client.execCommand resp:%v\n", resp.Success)
	// 如果结果为否
	if !resp.Success {
		// 如果有新的leader的address和port，更新
		leaderID := resp.GetLeaderID()
		if leaderID == "" {
			return false, errors.New("resp: leaderID == \"\"")
		}
		if leaderID != string(c.leader.ID) {
			c.leader.ID = raftlib.ServerID(leaderID)
			c.leader.Address = raftlib.ServerAddress(resp.GetLeaderAddress())
			c.leader.Port = raftlib.ServerPort(resp.GetLeaderPort())
			// 重新请求
			c.runRpcClient(true)
			return c.execCommand(logType, command)
		}
	}
	return resp.Success, nil
}

func main() {
	config.DEBUG = true
	cbFunc := func(b []byte) error {
		fmt.Println("fmt.Println:", string(b))
		return nil
	}
	// 创建节点
	config.Localserver = config.Servers[1]
	node := raftlib.CreateRaft(config)
	// 注册日志提交回调函数
	node.RegisterCallBackFunc(cbFunc)
	go node.Run()
	<-time.After(time.Second * 1)
	// 创建client
	config.Localserver = config.Servers[0]
	localServer = config.Servers[0]
	client := createClient(localServer)
	<-time.After(time.Second * 2)
	resp, _ := client.execCommand(raftlib.LogCommand, []byte("set testleft left"))
	fmt.Println("TestClient: recp:", resp)
	<-time.After(time.Second * 1)
}
