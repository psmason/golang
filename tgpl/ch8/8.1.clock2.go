package main

import (
	"io"
	"flag"
	"log"
	"net"
	"strconv"
	"time"
)

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05"))
		if err != nil {
			return // client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	port := flag.Int("port", 8000, "an integer corresponding to a port")
	flag.Parse()
	
	listener, err := net.Listen("tcp", "localhost:" + strconv.Itoa(*port))
	if err != nil {
		log.Fatal(err)	
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
