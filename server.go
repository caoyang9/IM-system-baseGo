package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex // 读写锁

	// 服务器广播消息的channel
	Msg chan string
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	// 创建server对象
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Msg:       make(chan string),
	}
	return server
}

// 监听服务器消息负责广播的goroutine
func (this *Server) ListenMsg() {
	for {
		msg := <-this.Msg
		this.mapLock.RLock()
		// 将msg发送给全部在线的user
		for _, client := range this.OnlineMap {
			client.C <- msg
		}
		this.mapLock.RUnlock()
	}
}

func (this *Server) Handler(conn net.Conn) {
	// 当前连接的业务
	//fmt.Println("连接建立成功...")
	user := NewUser(conn)

	// 用户上线，将用户加入到OnlineMap中
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 广播当前用户上线的消息
	this.Broadcast(user, "已上线")

	// 当前handler阻塞
	select {}
}

// 广播当前用户上线的消息
func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	// 向服务器的msg中发送用户上线消息
	this.Msg <- sendMsg
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

	// 监听服务器msg，广播给全部在线用户
	go this.ListenMsg()

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
