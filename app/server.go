package main

import (
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

// command map
type Command string

const (
	PING Command = "PING"
	ECHO Command = "ECHO"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection")
			os.Exit(1)
		}

		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		_, err := conn.Read(buf)

		if err != nil {
			fmt.Println("Failed to read data")
			return
		}

		fmt.Print(string(buf))

		requestLines := strings.Split(string(buf), "\r\n")

		if len(requestLines) < 3 {
			fmt.Println("Invalid request")
			return
		}

		command := Command(strings.ToUpper(strings.Trim(requestLines[2], "$")))

		var response string

		switch command {
		case PING:
			response = "+PONG\r\n"
		case ECHO:
			if len(requestLines) < 5 {
				fmt.Println("Invalid request")
				return
			}
			message := requestLines[4]
			response = "$" + fmt.Sprint(len(message)) + "\r\n" + message + "\r\n"

		default:
			response = "-ERR unknown command\r\n"
		}

		_, err = conn.Write([]byte(response))

		if err != nil {
			fmt.Println("Failed to write data")
			return
		}
	}

}
