package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
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
	cwdCommand    = "CWD"
	getCommand    = "RETR"
	quitCommand   = "QUIT"
	
	greeting            = "220 hello!"
	goodbye             = "221 goodbye"
	authenticated       = "230 ok we're good!"
	system              = "215 pile_of_garbage"
	portResponse        = "200 PORT successful"
	listOpen            = "150 opening for LIST"
	listComplete        = "226 LIST completed"
	typeResponse        = "200 Type set to A"
	unsupportedResponse = "500 unknown command"
	changedDirectory    = "250 Okay"
	noSuchDirectory     = "550 no such directory"
	noSuchFile          = "551 no such file to retrieve"
	fileWritten         = "226 file successfully written"
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

func formatLineFeeds(in string) string {
	return strings.Replace(in, "\n", "\r\n", -1)
}

func processList(commandConnection,
	dataConnection net.Conn,
	currentDirectory string) {
	// RFC 959, section 3.4. Transmission modes
	// http://stackoverflow.com/questions/37187986/bare-linefeeds-received-in-ascii-mode-warning-when-listing-directory-on-my-ftp
	// see above for carriage return usage
	log.Printf("LIST for current directory %s\n", currentDirectory)
	writeToConnection(commandConnection, listOpen+"\n")
	out, _ := exec.Command("ls", "-l", currentDirectory).Output()
	data := formatLineFeeds(string(out))
	writeToConnection(dataConnection, data+"\r")
	dataConnection.Close()
	writeToConnection(commandConnection, listComplete+"\n")
}

func processCwd(commandConnection net.Conn,
	currentDirectory *string,
	dir string) {

	log.Printf("Current working directory is %s\n", *currentDirectory)
	path := path.Clean(path.Join(*currentDirectory, strings.TrimSpace(dir)))
	log.Printf("Changing working directory to %s\n", path)

	fd, err := os.Stat(path)
	if err != nil {
		log.Printf("Failed to open Stat for %s, %s\n", path, err)
		writeToConnection(commandConnection, noSuchDirectory+"\n")
		return
	}

	if !fd.IsDir() {
		log.Printf("Path is not a directory: %s\n", path)
		writeToConnection(commandConnection, noSuchDirectory+"\n")
		return
	}

	*currentDirectory = path
	writeToConnection(commandConnection, changedDirectory+"\n")
}

func processGet(commandConnection,
	dataConnection net.Conn,
	currentDirectory string,
	file string) {

	file = path.Join(currentDirectory, strings.TrimSpace(file))
	log.Printf("Retrieving file %s\n", file)

	fi, err := os.Stat(file)
	if err != nil {
		log.Printf("Failed to open Stat for %s, %s\n", file, err)
		writeToConnection(commandConnection, noSuchFile+"\n")
		return
	}

	if fi.IsDir() {
		log.Printf("File is not a regular file: %s\n", file)
		writeToConnection(commandConnection, noSuchFile+"\n")
		return
	}

	writeToConnection(commandConnection, listOpen+"\n")
	
	fd, err := os.Open(file)
	defer fd.Close()
    if err != nil {
		log.Printf("Failed to open file %s\n", file)
		writeToConnection(commandConnection, noSuchFile+"\n")
		return
    }

	log.Printf("Copying file...")
	written, err := io.Copy(dataConnection, fd)
	dataConnection.Close()
	if err != nil {
		log.Printf("Failed to copy over file\n")
	}
	log.Printf("Copied over %d\n", written)
	writeToConnection(commandConnection, fileWritten+"\n")
}

func processQuit(c net.Conn) {
	writeToConnection(c, goodbye+"\n")
}

func processUnknown(c net.Conn) {
	writeToConnection(c, unsupportedResponse+"\n")
}

func processCommand(commandConnection net.Conn,
	dataConnection *net.Conn,
	currentDirectory *string,
	commandData string) {
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
		processList(commandConnection, *dataConnection, *currentDirectory)
	case typeCommand:
		processType(commandConnection)
	case cwdCommand:
		processCwd(commandConnection, currentDirectory, remainder)
	case getCommand:
		processGet(commandConnection,
			*dataConnection,
			*currentDirectory,
			remainder)
	case quitCommand:
		processQuit(commandConnection)
	default:
		processUnknown(commandConnection)
	}
}

func handleCommandConnection(commandConnection net.Conn) {
	defer commandConnection.Close()

	writeToConnection(commandConnection, greeting+"\n")
	var dataConnection net.Conn
	currentDirectory := "/home/patrick/dev/golang/tgpl/ch8"

	commandBuffer := make([]byte, commandBufferSize)
	bufferPos := 0
	for {
		if n, err := commandConnection.Read(commandBuffer); err != nil {
			log.Printf("Command connection closed: %s\n", err)
			return
		} else {
			bufferPos = n
		}

		processCommand(commandConnection,
			&dataConnection,
			&currentDirectory,
			string(commandBuffer[:bufferPos]))
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
