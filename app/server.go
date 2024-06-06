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

		for {
			buf := make([]byte, 1024)
			_, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading data: ", err.Error())
				break
			}
			splitStr := strings.Split((string(buf)), CLRF)
			if splitStr[0] == "GET / HTTP/1.1" {
				conn.Write([]byte(VERSION + " 200 OK" + CLRF + CLRF))
			} else {
				conn.Write([]byte(VERSION + " 404 Not Found" + CLRF + CLRF))
			}
		}
		conn.Close()
	}
}
