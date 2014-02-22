package muduo

import (
	// "log"
	"net"
)

func Check(v bool, msg string) {
	if !v {
		panic(msg)
	}
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func ListenTcpOrDie(listenAddr string) net.Listener {
	listener, err := net.Listen("tcp", listenAddr)
	PanicOnError(err)
	return listener
}
