package main

import (
	"fmt"
	"net"
	"os"
    "strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
    defer l.Close()
    fmt.Println("Server listnening on port 4221...")


    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting connection: ", err.Error())
            os.Exit(1)
        }
        if isValidRequest(conn) {
            conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
        } else {
            conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
        }
        conn.Close()
    }
}

func isValidRequest(conn net.Conn) bool {
    b := make([]byte, 1024)
    conn.Read(b)
    request := string(b)
    fmt.Println("Recieved request: ", request)
    parts := strings.Split(request, "\r\n")
    requestLine := parts[0]
    rl := strings.Split(requestLine, " ")
    path := rl[1]
    return path == "/"
}
