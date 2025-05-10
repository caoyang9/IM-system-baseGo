package main

import "net"

// 客户端用户结构体
type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server // User与Server绑定
}

// 用户上线
func (this *User) Online() {
	// 用户上线，将用户加入到OnlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 广播当前用户上线的消息
	this.server.Broadcast(this, "已上线")
}

// 用户下线
func (this *User) Offline() {
	// 用户下线，将用户从OnlineMap中剔除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 广播当前用户下线的消息
	this.server.Broadcast(this, "离线")
}

// 用户消息处理
func (this *User) DoMessage(msg string) {
	this.server.Broadcast(this, msg)
}

// 创建一个客户端用户实体
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}
	// 监听自身的channel，向conn中发消息
	go user.ListenMsg()

	return user
}

// 监听当前User中的channel
func (this *User) ListenMsg() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
