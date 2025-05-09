package main

import "net"

// 客户端用户结构体
type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// 创建一个客户端用户实体
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
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
