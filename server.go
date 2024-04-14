package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

// Server server类型 包含ip与port
type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex //锁

	//消息广播的channel
	Message chan string
}

// NewServer 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// ListenMsg 监听Message广播消息channel的goroutine,一旦有消息就发送给全部的在线User
func (s *Server) ListenMsg() {
	for {
		msg := <-s.Message

		onlineName := msg[strings.Index(msg, "]")+1 : strings.Index(msg, ":")]

		//将消息发送给除该上线用户外全部的在线User
		s.mapLock.Lock()
		//遍历OnlineMap,获取value
		for _, cli := range s.OnlineMap {
			if cli.Name == onlineName {
				continue
			}
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}

}

// BroadCast 广播消息
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr.String() + "]" + user.Name + ": " + msg

	//发送到Server Message channel
	s.Message <- sendMsg
}

// Handler 处理连接业务
func (s *Server) Handler(conn net.Conn) {
	//...当前连接的业务
	//fmt.Println("连接建立成功")

	//创建User
	user := NewUser(conn)

	//用户上线, 将用户加入到onlineMap中 (加锁)
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	//广播当前用户上线消息
	s.BroadCast(user, "已上线")

	//当前handler阻塞
	select {}
}

// Start 启动服务器的接口
func (s *Server) Start() {

	//socket listen
	//Sprintf:  format the variable as a string in the specified format
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))

	//judge err
	if errPrint(err, "net.listen err:") {
		return
	}

	//close listen socket
	defer listener.Close()

	//启动监听Message的goroutine
	go s.ListenMsg()

	for {
		//accept
		conn, err := listener.Accept()
		if errPrint(err, "listener accept err") {
			continue
		}
		//开一个协程do handler
		go s.Handler(conn)
	}
}

// ErrPrint judge and print err
func errPrint(err error, output string) bool {
	if err != nil {
		fmt.Println(output, err)
		return true
	}
	return false
}
