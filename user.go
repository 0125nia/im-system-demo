package main

import (
	"net"
	"strings"
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

// Online 用户上线的业务
func (u *User) Online(server *Server) {

	//用户上线,将用户加入到onlineMap中
	server.mapLock.Lock()
	server.OnlineMap[u.Name] = u
	server.mapLock.Unlock()

	//server广播用户上线消息
	server.BroadCast(u, "已上线")
}

// Offline 用户下线的业务
func (u *User) Offline(server *Server) {
	//用户下线,将用户从onlineMap中删除
	server.mapLock.Lock()
	delete(server.OnlineMap, u.Name)
	server.mapLock.Unlock()

	//广播当前用户下线消息
	server.BroadCast(u, "下线")
}

// SendMsg 给当前User对应的客户端发送消息
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

// DoMsg 用户处理消息的业务
func (u *User) DoMsg(server *Server, msg string) {

	if msg == "who" {
		server.mapLock.Lock()
		for _, user := range server.OnlineMap {
			onlineMsg := "[" + user.Name + "]:" + "在线...\n"
			u.SendMsg(onlineMsg)
		}
		server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//消息格式 rename|zhangSan
		newName := strings.Split(msg, "|")[1]

		//判断name是否存在
		_, ok := server.OnlineMap[newName]
		if ok {
			u.SendMsg("当前用户名被使用\n")
		} else {
			server.mapLock.Lock()
			delete(server.OnlineMap, u.Name)
			server.OnlineMap[newName] = u
			server.mapLock.Unlock()

			u.Name = newName
			u.SendMsg("您已更新用户名" + u.Name + "\n")
		}

	} else {
		server.BroadCast(u, msg)
	}
}

// ListenMsg 监听当前User channel的方法,一旦有消息,直接发送给对端客户端
func (u *User) ListenMsg() {
	for {
		msg := <-u.C

		u.conn.Write([]byte(msg + "\n"))
	}
}
