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
        b := make([]byte, 1024)
        n, err := conn.Read(b)
        if err != nil {
            fmt.Println("Error reading connection: ", err.Error())
            os.Exit(1)
        }
        fmt.Println("Read byte size: ", n)

        request := string(b)
        if strings.HasPrefix(request, "GET / HTTP/1.1\r\n") {
            conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
        } else {
            conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
        }
        conn.Close()
    }
}
