package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const VERSION = "HTTP/1.1"
const CLRF = "\r\n"

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	// Implement this function
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading data: ", err.Error())
			break
		}
		request := string(buf[:n])
		splitStr := strings.Split(request, CLRF)
		request_url_tokens := strings.Split(splitStr[0], " ")
		method := request_url_tokens[0]
		url := request_url_tokens[1]

		if method == "GET" {
			if url == "/" {
				response := VERSION + " 200 OK" + CLRF + CLRF
				conn.Write([]byte(response))
			} else if strings.Split(url, "/")[1] == "echo" {
				message := strings.Split(url, "/")[2]
				status := VERSION + " 200 OK" + CLRF
				headers := "Content-Type: text/plain" + CLRF + CLRF + fmt.Sprintf("Content-Length: %d", len(message)) + CLRF
				response := fmt.Sprintf("%s%s%s", status, headers, message)
				conn.Write([]byte(response))
			} else {
				response := VERSION + " 404 Not Found" + CLRF + CLRF
				conn.Write([]byte(response))
			}
		} else {
			response := VERSION + " 405 Method Not Allowed" + CLRF + "Content-Type: text/html" + CLRF + CLRF + "<h1>405 Method Not Allowed</h1>"
			conn.Write([]byte(response))
		}
	}
}
