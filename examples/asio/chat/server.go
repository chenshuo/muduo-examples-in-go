package chat

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
	// "sync/atomic"

	"github.com/chenshuo/muduo-examples-in-go/muduo"
)

type ChatServer struct {
	listener net.Listener
	nextConn int64
	mu       sync.Mutex
	conns    map[int64]net.Conn
}

func NewChatServer(listenAddr string) *ChatServer {
	server := new(ChatServer)
	server.listener = muduo.ListenTcpOrDie(listenAddr)
	server.conns = make(map[int64]net.Conn)
	return server
}

func readMessage(conn net.Conn, buf []byte) (n int, err error) {
	var length int32
	// FIXME: use bufio to save syscall
	err = binary.Read(conn, binary.BigEndian, &length)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if length > int32(len(buf)) || length < 0 {
		log.Println("Invalid length in header", length)
		return 0, errors.New("invalid length")
	}
	if length > 0 {
		var n int
		n, err = io.ReadFull(conn, buf[4:4+length])
		if err != nil {
			log.Println("error when reading payload", err)
			return 0, err
		}
		if int32(n) != length {
			log.Println("short reads", n, "of", length)
			return 0, errors.New("short read")
		}
	}
	return int(length), nil
}

func (s *ChatServer) ServeConn(conn net.Conn) {
	defer conn.Close()
	s.mu.Lock()
	connId := s.nextConn
	s.nextConn++
	s.conns[connId] = conn
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		delete(s.conns, connId)
		s.mu.Unlock()
	}()

	buf := make([]byte, 65536)
	for {
		length, err := readMessage(conn, buf)
		if err != nil {
			break
		}

		binary.BigEndian.PutUint32(buf, uint32(length))
		log.Println("Sending")
		conn.Write(buf)
	}
}

func (s *ChatServer) Run() {
	ticks := time.Tick(time.Second * 5)
	go func() {
		for _ = range ticks {
			s.mu.Lock()
			n := len(s.conns)
			s.mu.Unlock()
			fmt.Println(n)
		}
	}()
	muduo.ServeTcp(s.listener, s, "chat")
}
