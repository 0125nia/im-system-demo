package main

import . "Golang-IM-System/server"

func main() {
	//测试Server
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}
