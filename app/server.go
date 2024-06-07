package main

import (
	"flag"
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

	directory := flag.String("directory", ".", "Directory to serve files from")

	flag.Parse()

	if *directory != "" {
		os.Chdir(*directory)
	}

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
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
	}
	// Getting the request
	request := string(buf[:n])

	// getting the url of the request, reading until the first CLRF
	request_url := request[:(strings.Index(request, CLRF))]

	// getting the body and headers of the request
	request_body_and_header := request[(strings.Index(request, CLRF) + 2):]

	// splitting the body and headers using CLRF+CLRF
	request_body_and_header_tokens := strings.Split(request_body_and_header, CLRF+CLRF)

	// splitting the headers using CLRF
	headers := strings.Split(request_body_and_header_tokens[0], CLRF)

	// getting the body of the request
	body := request_body_and_header_tokens[1]

	// splitting the url to get the method and the url
	request_url_tokens := strings.Split(request_url, " ")

	// getting the method and the url
	method := request_url_tokens[0]
	url := request_url_tokens[1]

	// request headers
	request_headers := make(map[string]string)
	for _, header := range headers {
		if strings.TrimSpace(header) != "" {
			header_tokens := strings.Split(header, ": ")
			request_headers[header_tokens[0]] = header_tokens[1]
		}
	}

	if method == "GET" {
		if url == "/" {
			response := VERSION + " 200 OK" + CLRF + CLRF
			conn.Write([]byte(response))
		} else if strings.Split(url, "/")[1] == "echo" {
			message := strings.Split(url, "/")[2]
			status := VERSION + " 200 OK" + CLRF
			headers := "Content-Type: text/plain" + CLRF + fmt.Sprintf("Content-Length: %d", len(message)) + CLRF
			value, ok := request_headers["Accept-Encoding"]
			if ok {
				if value != "invalid-encoding" {
					headers += "Content-Encoding: " + value + CLRF
				}
			}
			headers += CLRF
			response := fmt.Sprintf("%s%s%s", status, headers, message)
			conn.Write([]byte(response))
		} else if strings.Split(url, "/")[1] == "user-agent" {
			user_agent := request_headers["User-Agent"]
			status := VERSION + " 200 OK" + CLRF
			headers := "Content-Type: text/plain" + CLRF + fmt.Sprintf("Content-Length: %d", len(user_agent)) + CLRF + CLRF
			response := fmt.Sprintf("%s%s%s", status, headers, user_agent)
			conn.Write([]byte(response))
		} else if strings.Split(url, "/")[1] == "files" {
			file_name := strings.Split(url, "/")[2]
			file, err := os.Open(file_name)
			if err != nil {
				response := VERSION + " 404 Not Found" + CLRF + CLRF
				conn.Write([]byte(response))
			} else {
				status := VERSION + " 200 OK" + CLRF
				headers := "Content-Type: application/octet-stream" + CLRF
				file_info, _ := file.Stat()
				file_size := file_info.Size()
				headers += fmt.Sprintf("Content-Length: %d", file_size) + CLRF + CLRF
				buf := make([]byte, file_size)
				file.Read(buf)
				response := fmt.Sprintf("%s%s%s", status, headers, buf)
				conn.Write([]byte(response))
			}
		} else {
			response := VERSION + " 404 Not Found" + CLRF + CLRF
			conn.Write([]byte(response))
		}
	} else if method == "POST" {
		if strings.Split(url, "/")[1] == "files" {
			file_name := strings.Split(url, "/")[2]
			file, err := os.Create(file_name)
			if err != nil {
				response := VERSION + "500 Internal Server Error" + CLRF + CLRF
				conn.Write([]byte(response))
			} else {
				file.Write([]byte(body))
				response := VERSION + " 201 Created" + CLRF + CLRF
				conn.Write([]byte(response))
			}
		}

	} else {
		response := VERSION + " 405 Method Not Allowed" + CLRF + CLRF
		conn.Write([]byte(response))
	}

}
