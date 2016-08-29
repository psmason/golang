package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	commandBufferSize = 1024

	userCommand   = "USER"
	systemCommand = "SYST"
	portCommand   = "PORT"
	listCommand   = "LIST"
	typeCommand   = "TYPE"

	greeting            = "220 hello!"
	authenticated       = "230 ok we're good!"
	system              = "215 pile_of_garbage"
	portResponse        = "200 PORT successful"
	listOpen            = "150 opening for LIST"
	listComplete        = "226 LIST completed"
	typeResponse        = "200 Type set to A"
	unsupportedResponse = "500 unknown command"
)

func writeToConnection(c net.Conn, data string) {
	log.Printf("Writing to connection: %s\n", data)
	_, err := io.WriteString(c, data)
	if err != nil {
		log.Fatal("Failed to write to connection")
	}
}

func processUser(c net.Conn) {
	// notice this doesn't really do anything secure
	writeToConnection(c, authenticated+"\n")
}

func processSyst(c net.Conn) {
	writeToConnection(c, system+"\n")
}

func processType(c net.Conn) {
	writeToConnection(c, typeResponse+"\n")
}

func processPort(commandConnection net.Conn, dataConnection *net.Conn, destinationString string) {
	log.Printf("Handling PORT for %s\n", destinationString)

	tokens := strings.Split(destinationString, ",")
	ip := strings.Join(tokens[:len(tokens)-2], ".")

	port1, _ := strconv.Atoi(tokens[len(tokens)-2])
	port2, _ := strconv.Atoi(strings.TrimSpace(tokens[len(tokens)-1]))
	destination := ip + ":" + strconv.Itoa(port1*256+port2)
	log.Printf("Connecting to destination %s\n", destination)

	tmpConnection, err := net.Dial("tcp", destination)
	if err != nil {
		log.Fatal(err)
	}
	*dataConnection = tmpConnection

	writeToConnection(commandConnection, portResponse+"\n")
}

func processList(commandConnection, dataConnection net.Conn) {
	// RFC 959, section 3.4. Transmission modes
	// http://stackoverflow.com/questions/37187986/bare-linefeeds-received-in-ascii-mode-warning-when-listing-directory-on-my-ftp
	// see above for carriage return usage
	writeToConnection(commandConnection, listOpen+"\n")
	out, _ := exec.Command("ls").Output()
	data := strings.Replace(string(out), "\n", "\r\n", -1)
	writeToConnection(dataConnection, data+"\r")
	dataConnection.Close()
	writeToConnection(commandConnection, listComplete+"\n")
}

func processUnknown(c net.Conn) {
	writeToConnection(c, unsupportedResponse+"\n")
}

func processCommand(commandConnection net.Conn, dataConnection *net.Conn, commandData string) {
	log.Printf("Processing command data %s\n", commandData)
	tokens := strings.Split(commandData, " ")
	command := strings.TrimSpace(tokens[0])
	remainder := strings.Join(tokens[1:], "")
	log.Printf("Processing command %s\n", command)
	switch command {
	case userCommand:
		processUser(commandConnection)
	case systemCommand:
		processSyst(commandConnection)
	case portCommand:
		processPort(commandConnection, dataConnection, remainder)
	case listCommand:
		processList(commandConnection, *dataConnection)
	case typeCommand:
		processType(commandConnection)
	default:
		processUnknown(commandConnection)
	}
}

func handleCommandConnection(commandConnection net.Conn) {
	defer commandConnection.Close()

	writeToConnection(commandConnection, greeting+"\n")
	var dataConnection net.Conn

	commandBuffer := make([]byte, commandBufferSize)
	bufferPos := 0
	for {
		if n, err := commandConnection.Read(commandBuffer); err != nil {
			log.Fatal(err)
		} else {
			bufferPos = n
		}

		processCommand(commandConnection, &dataConnection, string(commandBuffer[:bufferPos]))
	}
}

func commandListener() {
	listener, err := net.Listen("tcp", "localhost:8010")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleCommandConnection(conn)
	}
}

func main() {
	log.SetOutput(os.Stdout)
	go commandListener()

	forever := make(chan bool)
	<-forever
}
