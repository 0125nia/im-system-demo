package main

import (
	"Golang-IM-System/util"
	"flag"
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

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认为127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认为8888)")
}

// main 开启客户端
func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)

	if client == nil {
		fmt.Println(">>>>>>> 连接服务器失败 >>>>>>>")
		return
	}

	fmt.Println(">>>>>>> 连接服务器成功 >>>>>>>")

	//启动Client的业务
	select {}
}
