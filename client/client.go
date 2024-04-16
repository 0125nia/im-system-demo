package main

import (
	"Golang-IM-System/util"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

// Client 客户端对象
type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //当前client的模式
}

// NewClient 创建客户端对象
func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
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

// DealResponse 处理server回应的消息，直接显示到标准输出即可
func (c *Client) DealResponse() {
	//一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, c.conn)
}

// menu 模式选择目录
func (c *Client) menu() bool {
	var f int
	fmt.Println(">>>>>>> 1.公聊 <<<<<<<")
	fmt.Println(">>>>>>> 2.私聊 <<<<<<<")
	fmt.Println(">>>>>>> 3.改名 <<<<<<<")
	fmt.Println(">>>>>>> 0.退出 <<<<<<<")

	fmt.Scanln(&f)

	if f >= 0 && f <= 3 {
		c.flag = f
		return true
	} else {
		fmt.Println(">>>>>>> 请输入合法范围内的数字 <<<<<<<")
		return false
	}
}

// SelectUsers 查询在线用户
func (c *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if util.ErrPrint(err, "conn Write") {
		return
	}
}

// PrivateChat 私聊模式
func (c *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	c.SelectUsers()
	fmt.Println(">>>>>>> 请输入聊天对象[用户名]，exit退出：")
	fmt.Scanln(&remoteName)
	for remoteName != "exit" {
		fmt.Println(">>>>>>> 请输入消息内容，exit退出：")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := c.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>>>> 请输入聊天内容，exit退出")
			fmt.Scanln(&chatMsg)
		}

		c.SelectUsers()
		fmt.Println(">>>>>>> 请输入聊天对象[用户名]，exit退出：")
		fmt.Scanln(&remoteName)
	}

}

// PublicChat 公聊模式
func (c *Client) PublicChat() {
	//提示用户输入消息
	var chatMsg string

	fmt.Println(">>>>>>> 请输入聊天内容，exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>>>> 请输入聊天内容，exit退出")
		fmt.Scanln(&chatMsg)
	}
}

// UpdateName 更改用户名
func (c *Client) UpdateName() bool {
	fmt.Println(">>>>>>> 请输入用户名：")
	fmt.Scanln(&c.Name)

	sendMsg := "rename|" + c.Name + "\n"

	_, err := c.conn.Write([]byte(sendMsg))
	if util.ErrPrint(err, "conn.Write") {
		return false
	}
	return true
}

// Run 执行Client业务
func (c *Client) Run() {
	for c.flag != 0 {
		for !c.menu() {

		}

		//根据不同模式处理不同业务
		switch c.flag {
		case 1:
			//公聊模式
			fmt.Println("选择公聊模式...")
			c.PublicChat()
			break
		case 2:
			//私聊模式
			fmt.Println("选择私聊模式...")
			c.PrivateChat()
			break
		case 3:
			//改名
			fmt.Println("选择改名模式...")
			c.UpdateName()
			break
		}
	}
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

	//单独开启一个goroutine去处理server的回执消息
	go client.DealResponse()

	fmt.Println(">>>>>>> 连接服务器成功 >>>>>>>")

	//启动Client的业务
	client.Run()
}
