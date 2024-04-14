package main

import (
	"Golang-IM-System/util"
	"fmt"
	"net"
)

// Client 客户端对象
type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

// NewClient 创建客户端对象
func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	//连接Server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if util.ErrPrint(err, "net.Dial") {
		return nil
	}
	client.conn = conn
	//返回对象
	return client
}

// main 开启客户端
func main() {
	client := NewClient("127.0.0.1", 8888)

	if client == nil {
		fmt.Println(">>>>>>> 连接服务器失败 >>>>>>>")
		return
	}

	fmt.Println(">>>>>>> 连接服务器成功 >>>>>>>")

	//启动Client的业务
	select {}
}
