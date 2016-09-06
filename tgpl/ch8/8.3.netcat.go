package main

import (
	"io"
	"log"
	"net"
	"os"
)

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetOutput(os.Stdout)
	tcpAddress,_ := net.ResolveTCPAddr("tcp", "localhost:8000")
	conn, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	go func() {
		io.Copy(os.Stdout, conn)
		log.Println("done")
		done <- true
	}()
	mustCopy(conn, os.Stdin)
	conn.CloseWrite()
	<-done
}


