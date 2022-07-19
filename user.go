package main

import (
	"net"
)

type User struct {
	Name    string
	Address string
	Channel chan string
	conn    net.Conn
	server  *Server
}

// 构造方法
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:    userAddr,
		Address: userAddr,
		Channel: make(chan string),
		conn:    conn,
		server:  server,
	}

	// 启动监听
	go user.ListenMessage()
	return user
}

// 用户上线业务
func (this *User) Online() {

	// 用户上线,将用户加入到onlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 广播用户上线
	this.server.BroadCast(this, "login")
}

// 用户下线业务
func (this *User) Offline() {
	// 用户下线,将用户从onlineMap中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 广播用户上线
	this.server.BroadCast(this, "logout")
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
}

// 监听channel
func (this *User) ListenMessage() {
	for true {
		msg := <-this.Channel
		this.conn.Write([]byte(msg + "\n"))
	}
}
