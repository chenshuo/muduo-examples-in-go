package main

import (
	"github.com/chenshuo/muduo-examples-in-go/examples/asio/chat"
)

func main() {
	server := chat.NewChatServer(":3399")
	server.Run()
}
