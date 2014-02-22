package simple

import (
	"bytes"
	"io"
	"log"
	"net"
)

type ChargenServer struct {
	listener net.Listener
	// TODO: printThroughput()
}

var message []byte
var repeatReader *RepeatReader

func NewChargenServer(listenAddr string) *ChargenServer {
	server := new(ChargenServer)
	server.listener = listenTcpOrDie(listenAddr)
	return server
}

func init() {
	var buf bytes.Buffer
	for i := 33; i < 127; i++ {
		buf.WriteByte(byte(i))
	}
	l := buf.Len()
	n, _ := buf.Write(buf.Bytes())
	if l != n || buf.Len() != 2*l {
		panic("buf error")
	}
	var msg bytes.Buffer
	for i := 0; i < 127-33; i++ {
		msg.Write(buf.Bytes()[i : i+72])
		msg.WriteByte('\n')
	}
	message = msg.Bytes()
	if len(message) != 94*73 {
		panic("short message")
	}
	repeatReader = NewRepeatReader(message)
}

func chargen(c net.Conn) {
	defer c.Close()
	var total int64
	for {
		n, err := c.Write(message)
		total += int64(n)
		if err != nil {
			log.Println("chargen:", err.Error())
			break
		}
		if n != len(message) {
			log.Println("chargen: short write", n, "out of", len(message))
		}
	}
	log.Println("total", total)
	printConn(c, "chargen", "DOWN")
}

/////////////////// RepeatReader

type RepeatReader struct {
	message []byte
}

func NewRepeatReader(m []byte) *RepeatReader {
	r := new(RepeatReader)
	r.message = m
	return r
}

func (r *RepeatReader) Read(p []byte) (n int, err error) {
	if len(p) < len(r.message) {
		panic("Short read")
	}
	copy(p, r.message)
	return len(r.message), nil
}

/////////////////// RepeatReader

func chargenAdv(c net.Conn) {
	defer c.Close()
	total, err := io.Copy(c, repeatReader)
	if err != nil {
		log.Println("chargen:", err.Error())
	}
	log.Println("total", total)
	printConn(c, "chargen", "DOWN")
}

func (s *ChargenServer) Serve() {
	serveTcp(s.listener, chargenAdv, "chargen")
}
