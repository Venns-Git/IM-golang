package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip string
	Port int
}
// 创建一个server的接口
func NewServer(ip string,port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
	}
	return server
}
// 业务接口
func (this *Server) Handler(conn net.Conn)  {
	fmt.Println("连接建立成功")
}

// 启动服务的接口
func (this *Server) Start() {
	// socket listen

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err: ",err)
	}

	// close listen socket
	defer listener.Close()

	for true {
		// accept
		conn, err := listener.Accept()
		if err != nil{
			fmt.Println("listener accept err: ",err)
			continue
		}
		
		// do handler
		go this.Handler(conn)
	}
}
