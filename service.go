package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播
	Message chan string
}

// 监听message广播消息的channel
func (this *Server) ListenMessager() {
	for true {
		msg := <-this.Message

		// 将msg发送给全部的在线User
		this.mapLock.Lock()
		for _, client := range this.OnlineMap {
			client.Channel <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Address + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 业务接口
func (this *Server) Handler(conn net.Conn) {

	user := NewUser(conn, this)

	user.Online()

	// 监听用户是否活跃的channel
	isLive := make(chan bool)

	// 接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for true {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err: ", err)
				return
			}

			// 提取用户消息
			msg := string(buf[:n-1])

			// 用户针对msg进行消息处理
			user.DoMessage(msg)

			// 用户的任意消息,代表当前用户是一个活跃的
			isLive <- true
		}
	}()

	// 当前handler阻塞
	for {
		select {
		case <-isLive:
			// 当前用户是活跃的, 应该重置定时器
			// 不做任何事情,为了激活select,更新下面的定时器
		case <-time.After(time.Second * 10):
			// 已经超时
			// 将当前的客户端user强制关闭
			user.SendMsg("Timeout connection close")

			// 销毁资源
			close(user.Channel)

			// 关闭连接
			conn.Close()

			// 退出当前handler
			return
		}
	}
}

// 启动服务的接口
func (this *Server) Start() {
	// socket listen

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
	}

	// close listen socket
	defer listener.Close()

	// 启动监听message
	go this.ListenMessager()

	for true {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err: ", err)
			continue
		}

		// do handler
		go this.Handler(conn)
	}
}
func main() {
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}
