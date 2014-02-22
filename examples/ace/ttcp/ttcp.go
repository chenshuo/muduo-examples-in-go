package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/chenshuo/muduo-examples-in-go/muduo"
)

type options struct {
	port     int
	length   int
	number   int
	transmit bool
	receive  bool
	nodelay  bool
	host     string
}

var opt options

type SessionMessage struct {
	Number, Length int32
}

func init() {
	flag.IntVar(&opt.port, "p", 5001, "TCP port")
	flag.IntVar(&opt.length, "l", 65536, "Buffer length")
	flag.IntVar(&opt.number, "n", 8192, "Number of buffers")
	flag.BoolVar(&opt.receive, "r", false, "Receive")
	flag.StringVar(&opt.host, "t", "", "Transmit")
	if binary.Size(SessionMessage{}) != 8 {
		panic("packed struct")
	}
}

func transmit() {
	// t := conn.(*net.TCPConn)
	// t.SetNoDelay(false)
}

func receive() {
	listener := muduo.ListenTcpOrDie(fmt.Sprintf(":%d", opt.port))
	defer listener.Close()
	conn, err := listener.Accept()
	muduo.PanicOnError(err)
	defer conn.Close()

	// Read header
	var sessionMessage SessionMessage
	err = binary.Read(conn, binary.BigEndian, &sessionMessage)
	muduo.PanicOnError(err)

	fmt.Printf("receive buffer length = %d\nreceive number of buffers = %d\n",
		sessionMessage.Length, sessionMessage.Number)
	total_mb := float64(sessionMessage.Number) * float64(sessionMessage.Length) / 1024.0 / 1024.0
	fmt.Printf("%.3f MiB in total\n", total_mb)

	payload := make([]byte, sessionMessage.Length)
	start := time.Now()
	for i := 0; i < int(sessionMessage.Number); i++ {
		var length int32
		err = binary.Read(conn, binary.BigEndian, &length)
		muduo.PanicOnError(err)
		if length != sessionMessage.Length {
			panic("read length")
		}

		var n int
		n, err = io.ReadFull(conn, payload)
		muduo.PanicOnError(err)
		if n != len(payload) {
			panic("read payload data")
		}

		// ack
		err = binary.Write(conn, binary.BigEndian, &length)
		muduo.PanicOnError(err)
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("%.3f seconds\n%.3f MiB/s\n", elapsed, total_mb/elapsed)
}

func main() {
	flag.Parse()
	opt.transmit = opt.host != ""
	if opt.transmit == opt.receive {
		println("Either -r or -t must be specified.")
		return
	}

	if opt.transmit {
		println("transmit")
	} else if opt.receive {
		// println("receive")
		receive()
	} else {
		panic("unknown")
	}
}
