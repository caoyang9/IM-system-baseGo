package main

import (
	"net"
	"strings"
)

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

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 用户消息处理
func (this *User) DoMessage(msg string) {
	stringArr := strings.Split(msg, " ")
	if msg == "im online\n" {
		// 查询当前所有在线用户
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			// 写到当前用户的channel
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if stringArr[1] == "rename" {
		// 去掉新用户名后的空格
		newName := stringArr[2][:len(stringArr[2])-1]
		// 判断用户是否存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("该用户名已存在！\n")
		} else {
			this.server.mapLock.Lock()
			// 删除map中的原用户，新增新名称的用户
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMsg("更新用户名成功：" + this.Name + "\n")
		}
	} else {
		this.server.Broadcast(this, msg)
	}
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
