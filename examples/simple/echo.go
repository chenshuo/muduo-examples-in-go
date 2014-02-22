package simple

import (
	"io"
	"log"
	"net"
)

type EchoServer struct {
	listener net.Listener
}

func NewEchoServer(listenAddr string) *EchoServer {
	server := new(EchoServer)
	server.listener = listenTcpOrDie(listenAddr)
	return server
}

func echo(c net.Conn) {
	defer c.Close()
	total, err := io.Copy(c, c)
	if err != nil {
		log.Println("echo:", err.Error())
	}
	log.Println("total", total)
	printConn(c, "echo", "DOWN")
}

func (s *EchoServer) Serve() {
	serveTcp(s.listener, echo, "echo")
}
