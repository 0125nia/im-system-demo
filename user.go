package main

import (
	"net"
)

// User User类型 用户结构体
type User struct {
	Addr net.Addr
	Name string
	C    chan string
	conn net.Conn
}

// NewUser 创建一个User的API
func NewUser(conn net.Conn) *User {
	//获取客户端连接地址
	userAddr := conn.RemoteAddr()
	userName := userAddr.String()

	user := &User{
		Addr: userAddr,
		Name: userName,
		C:    make(chan string),
		conn: conn,
	}

	//启动监听当前user channel消息的goroutine
	go user.ListenMsg()

	return user
}

// ListenMsg 监听当前User channel的方法,一旦有消息,直接发送给对端客户端
func (u *User) ListenMsg() {
	for {
		msg := <-u.C

		u.conn.Write([]byte(msg + "\n"))
	}
}
