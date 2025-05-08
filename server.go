package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	// 创建server对象
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (this *Server) Handler(conn net.Conn) {
	// 当前连接的业务
	fmt.Println("连接建立成功...")
}

// 启动服务器的接口
func (this *Server) Start() {
	// socket listening
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	// 关闭套接字
	defer listener.Close()

	for {
		// 返回一个conn连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}

		// do handler 让go程去异步的执行任务
		go this.Handler(conn)
	}
}
