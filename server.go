package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Server server类型 包含ip与port
type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex //锁

	//消息广播的channel
	Message chan Message
}

// NewServer 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan Message),
	}
	return server
}

// ListenMsg 监听Message广播消息channel的goroutine,一旦有消息就发送给全部的在线User
func (s *Server) ListenMsg() {
	for {
		msg := <-s.Message
		name := msg.userName

		sendMsg := "[" + name + "]" + ": " + msg.msg

		//将消息发送给除该上线用户外全部的在线User
		s.mapLock.Lock()
		//遍历OnlineMap,获取value即User
		for _, cli := range s.OnlineMap {
			//若为该用户则不发送
			if cli.Name == name {
				continue
			}
			cli.C <- sendMsg
		}
		s.mapLock.Unlock()
	}

}

// BroadCast 广播消息
func (s *Server) BroadCast(user *User, msg string) {
	//发送到Server Message channel
	s.Message <- Message{user.Name, msg}
}

// Handler 处理连接业务
func (s *Server) Handler(conn net.Conn) {
	//...当前连接的业务
	//fmt.Println("连接建立成功")

	//创建User
	user := NewUser(conn)

	user.Online(s)

	//在用户的Handler() goroutine中添加用户活跃channel,一旦有消息即向该channel发送消息类似激活操作
	//长时间不发送消息将被认为超时,将被强制关闭用户连接
	//监听用户是否活跃的channel
	isLive := make(chan bool)

	//接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)

		for {
			//读取信息  n表示字节数组长度
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline(s)
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			//提取用户的消息(去除'\n')
			msg := string(buf[:n-1])

			//将得到的消息进行广播
			user.DoMsg(s, msg)

			//用户的任意消息,代表当前用户是一个活跃的状态
			isLive <- true
		}
	}()

	//当前handler阻塞
	for {
		//超时强制踢出功能
		select {
		case <-isLive:
			//当前用户是活跃的,应重置定时器
			//不做任何操作,由于case穿透会执行下一行判断操作,变相重置了定时器
		case <-time.After(time.Minute * 5): //在用户的goroutine中添加定时器功能,判断是否超时
			//执行此case代表已经超时
			//将当前的User连接强制关闭
			user.SendMsg("你因超时被踢了")

			//释放该用户所用的资源
			close(user.C)

			//关闭连接
			conn.Close()

			//退出当前Handler
			return //runtime.Goexit()
		}
	}

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
