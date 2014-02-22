package simple

import (
	"log"
	"net"
)

func printConn(c net.Conn, name, updown string) {
	log.Printf("%v: %v <- %v is %v\n",
		name, c.LocalAddr().String(), c.RemoteAddr().String(), updown)
}

type ServeTcpConn func(c net.Conn)

func serveTcp(l net.Listener, f ServeTcpConn, name string) {
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err == nil {
			// t := conn.(*net.TCPConn)
			// t.SetNoDelay(false)
			printConn(conn, name, "UP")
			go f(conn)
		} else {
			log.Println(name, err.Error())
		}
	}
}
