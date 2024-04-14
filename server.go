package main

import (
	"fmt"
	"net"
)

// Server server类型 包含ip与port
type Server struct {
	Ip   string
	Port int
}

// NewServer 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

// Handler 处理连接业务
func (this *Server) Handler(conn net.Conn) {
	//...当前连接的业务
	fmt.Println("连接建立成功")
}

// Start 启动服务器的接口
func (this *Server) Start() {

	//socket listen
	//Sprintf:  format the variable as a string in the specified format
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	//judge err
	if errPrint(err, "net.listen err:") {
		return
	}

	//close listen socket
	defer listener.Close()

	for {
		//accept
		conn, err := listener.Accept()
		if errPrint(err, "listener accept err") {
			continue
		}
		//开一个协程do handler
		go this.Handler(conn)
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
