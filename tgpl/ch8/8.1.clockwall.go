package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	ny := flag.String("NewYork", "localhost:8010", "an address for a NY clock")
	tk := flag.String("Tokyo", "localhost:8020", "an address for a Tokyo clock")
	ln := flag.String("London", "localhost:8030", "an address for a London clock")
	flag.Parse()
	nyConn, lnConn, tkConn := subscribe("NewYork", *ny), subscribe("London", *ln), subscribe("Tokyo", *tk)
	defer nyConn.Close()
	defer lnConn.Close()
	defer tkConn.Close()

	fmt.Printf("NewYork\t\tLondon\t\tTokyo\n");
	for {
		nyTime, lnTime, tkTime := readTime(nyConn), readTime(lnConn), readTime(tkConn)
		fmt.Printf("%s\t%s\t%s\n", nyTime, lnTime, tkTime)
	}
}

func subscribe(clockname, destination string) net.Conn {
	log.Print("subscribing to", destination)
	conn, err := net.Dial("tcp", destination)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func readTime(src io.Reader) string {
	timeBuffer := make([]byte, 8)
	bufferPos := 0
	if n, err := io.ReadFull(src, timeBuffer); err != nil {
		log.Fatal(err)
	} else {
		bufferPos = n
	}
	return string(timeBuffer[:bufferPos])
}
