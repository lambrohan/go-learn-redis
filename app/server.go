package main

import (
	"fmt"
	"strconv"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

// command map
type Command struct {
	Name CommandType
	Args []string
}
type CommandType string

const (
	INFO CommandType = "INFO"
	PING CommandType = "PING"
	ECHO CommandType = "ECHO"
	SET  CommandType = "SET"
	GET  CommandType = "GET"
)

var redisData = make(map[string]string)

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

		command, err := parseBuffer(buf)

		if err != nil {
			fmt.Println("Failed to read data")
			return
		}

		var response string

		switch command.Name {
		case INFO:
			response = "+OK\r\n"
		case PING:
			response = "+PONG\r\n"
		case ECHO:
			if len(command.Args) == 0 {
				response = "-ERR wrong number of arguments for 'echo' command\r\n"
			} else {
				response = formatResponse(command.Args[0])
			}
		case SET:
			if len(command.Args) != 2 {
				response = "-ERR wrong number of arguments for 'set' command\r\n"
			} else {
				redisData[command.Args[0]] = command.Args[1]
				response = "+OK\r\n"
			}
		case GET:
			if len(command.Args) != 1 {
				response = "-ERR wrong number of arguments for 'get' command\r\n"
			} else {
				value, ok := redisData[command.Args[0]]
				if ok {
					response = formatResponse(value)
				} else {
					response = "$-1\r\n"
				}

			}

		default:
			// response = "-ERR unknown command\r\n"
		}

		_, err = conn.Write([]byte(response))

		if err != nil {
			fmt.Println("Failed to write data")
			return
		}
	}

}

func formatResponse(response string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(response), response)
}

func parseBuffer(buf []byte) (Command, error) {
	requestLines := strings.Split(string(buf), "\r\n")

	if len(requestLines) < 3 {
		return Command{}, fmt.Errorf("Invalid request")
	}

	numArgs, err := strconv.Atoi(requestLines[0][1:])
	if err != nil {
		return Command{}, fmt.Errorf("Invalid request")
	}

	if len(requestLines) < numArgs*2+1 {
		return Command{}, fmt.Errorf("Invalid request")
	}

	command := Command{
		Name: CommandType(strings.ToUpper(requestLines[2])),
		Args: make([]string, numArgs-1),
	}

	for i := 0; i < numArgs-1; i++ {
		command.Args[i] = requestLines[i*2+4]
	}

	return command, nil
}
