# raftlib

## 介绍
- 毕设作业，实现简易raft库
- 使用go + grpc实现
- 主要是对Raft的三个子问题进行编码实现，功能不完善。

## 说明
只是个玩具，很多细节可能为考虑完善，比如说部分地方没有做超时限制，日志项储存用的简单字节数组缓存，没有实现日志压缩和集群成员变更的功能，没有特别完善的错误处理，无法安全关机等。

## 快速开始

### 配置文件
```go
type Config struct {
	// 是否启用Debug日志
	DEBUG bool

	// 设置超时时间
	// 心跳超时时间
	HeartbeatTimeout time.Duration
	// 选举超时时间
	ElectionTimeout time.Duration
	// RPC响应超时时间
	RpcTimeout time.Duration

	// Servers: 服务器元信息列表
	Servers []Server

    // 当前Raft节点
    localserver Server
}
```

### 配置集群成员
启动raft前，需要手动配置集群成员，即配置config中的Servers。

Server结构体定义如下：
```go
// Server: 服务器元信息
type Server struct {
	// Suffrage: 表示节点是否获得投票
	Suffrage ServerSuffrage
	// ID: 总是唯一确定的服务器标识码
	ID ServerID
	// Address: 服务器网络地址
	Address ServerAddress
	// Port: 服务器端口号
	Port ServerPort
}
```

其中Suffrage务必设置为Voter，因为时间原因没有做投票时机部分的编码实现，请见谅。

可以设置成如下集群配置：
```go
config := raftlib.CreateConfig()
config.Servers = []Server{
	{Voter, "0", "192.168.224.101", "5678"},
	{Voter, "1", "192.168.224.102", "5678"},
}
```

### 配置当前Raft节点
启动raft前，需要手动配置当前Raft节点，即配置config中的LocalServer，如：
```go
config.LocalServer = Server{Voter, "1", "192.168.224.102", "5678"}
```

### 创建Raft节点
使用raftlib.CreateRaft(config),可以返回一个raft句柄，其中config就是配置文件

### 注册日志提交回调函数
使用返回的raft句柄，调用registerCallBackFunc方法即可注册日志提交回调函数，其中方法定义如下：
```go
func (r *Raft) RegisterCallBackFunc(f func([]byte) error)
```

### 运行Raft
句柄为raft，执行以下方法
raft.Run()

### execCommandRPC请求
可以参考client中的execCommand方法

### example
client.go即为一个example，功能：client是分布式缓存的一个节点，使用了redis进行缓存存储。有需要的朋友可以对回调函数进行更换。

在main中多创建了一个节点，用于模拟raft集群运行。
