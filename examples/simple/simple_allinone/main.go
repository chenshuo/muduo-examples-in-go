package main

import (
	"github.com/chenshuo/muduo-examples-in-go/examples/simple"
)

func main() {
	ch := make(chan bool)

	chargenServer := simple.NewChargenServer(":2019")
	go chargenServer.Serve()

	discardServer := simple.NewDiscardServer(":2009")
	go discardServer.Serve()

	echoServer := simple.NewEchoServer(":2007")
	go echoServer.Serve()

	<-ch  // wait forever
}
