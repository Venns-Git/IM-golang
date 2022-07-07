package main

import "net"

type User struct {
	Name    string
	Address string
	Channel chan string
	conn    net.Conn
}

// 构造方法
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:    userAddr,
		Address: userAddr,
		Channel: make(chan string),
		conn:    conn,
	}

	// 启动监听
	go user.ListenMessage()
	return user
}

// 监听channel
func (this *User) ListenMessage() {
	for true {
		msg := <-this.Channel
		this.conn.Write([]byte(msg + "\n"))
	}
}
