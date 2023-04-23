package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	example "raftlib/example"
	"raftlib/pb"
	raftlib "raftlib/src"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var clientConfig *example.RaftConfig

// 客户端
type Client struct {
	ClientService pb.RpcServiceClient
	ClientConn    *grpc.ClientConn
	LeaderService pb.RpcServiceClient
	LeaderConn    *grpc.ClientConn
	Leader        raftlib.Server
	RedisClient   *redis.Client
}

func createClient(server raftlib.Server) *Client {
	client := &Client{
		Leader: server,
	}
	client.createRedisClient()
	// 创建RPC连接
	client.runRpcClient(true)
	client.runRpcClient(false)
	fmt.Printf("createClient finished.\n")
	return client
}

func (c *Client) createRedisClient() {
	// 通过配置文件读取Redis客户端参数
	redisClient := redis.NewClient(&redis.Options{
		Addr:        clientConfig.REDIS_ADDRESS,
		Password:    clientConfig.REDIS_PASSWORD,
		DB:          clientConfig.REDIS_DB,
		DialTimeout: clientConfig.REDIS_DIAL_TIMEOUT,
		ReadTimeout: clientConfig.REDIS_READ_TIMEOUT,
	})
	// 测试连接是否成功
	pong, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("createRedisClient(): pong: %v\n", pong)
	c.RedisClient = redisClient
}

// 启动rpc客户端
func (c *Client) runRpcClient(isLeader bool) error {
	fmt.Printf("runRpcClient, leaderAddress: %v, leaderPort: %v\n", c.Leader.Address, c.Leader.Port)
	// 创建客户端
	clientConn, err := grpc.Dial(string(c.Leader.Address)+":"+string(c.Leader.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 断线重连参数
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                2 * time.Second,
			Timeout:             500 * time.Millisecond,
			PermitWithoutStream: true}))
	if err != nil {
		log.Fatal(err)
		return err
	}
	// 创建RpcService
	clientService := pb.NewRpcServiceClient(clientConn)
	if isLeader {
		c.LeaderConn = clientConn
		c.LeaderService = clientService
	} else {
		c.ClientConn = clientConn
		c.ClientService = clientService
	}

	return nil
}

func (c *Client) execCommand(logType raftlib.LogType, command []byte) ([]byte, error) {
	fmt.Printf("Client.execCommand.\n")
	ver := raftlib.PROTO_VER_EXEC_COMMAND_REQUEST
	req := &pb.ExecCommandRequest{
		Ver:     &ver,
		LogType: uint32(logType),
		Command: command,
	}
	ctx := context.Background()
	resp, _ := c.LeaderService.ExecCommand(ctx, req)
	fmt.Printf("Client.execCommand resp:%v\n", resp.GetSuccess())
	// 如果结果为否
	if !resp.GetSuccess() {
		// 如果有新的leader的address和port，更新
		leaderID := resp.GetLeaderID()
		if leaderID == "" {
			return nil, errors.New("resp: leaderID == \"\"")
		}
		if leaderID != string(c.Leader.ID) {
			c.Leader.ID = raftlib.ServerID(leaderID)
			c.Leader.Address = raftlib.ServerAddress(resp.GetLeaderAddress())
			c.Leader.Port = raftlib.ServerPort(resp.GetLeaderPort())
			// 重新请求
			c.runRpcClient(true)
			return c.execCommand(logType, command)
		}
	}
	return resp.GetData(), nil
}

func main() {
	clientConfig = example.CreateConfig()
	client := createClient(clientConfig.LocalServer)
	for {
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		// 解析Redis命令
		commandSplited := strings.Split(command, " ")
		switch strings.ToLower(commandSplited[0]) {
		case "del":
			fallthrough
		case "set":
			resp, err := client.execCommand(raftlib.LogCommand, []byte(command))
			if err != nil {
				fmt.Println("exec command fail:", err)
			} else {
				fmt.Println(string(resp))
			}
		case "get":
			key := strings.Trim(commandSplited[1], " \n")
			result, err := client.RedisClient.Get(key).Result()
			if err != nil {
				fmt.Printf("Get error: %v\n", err)
			}
			fmt.Println(result)
		}
	}
}
