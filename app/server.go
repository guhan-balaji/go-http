package main

import (
        "fmt"
        "net"
        "os"
        "strings"
)

type RequestLine struct {
        verb, target, protocol string
}

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
                _, err = conn.Read(b)
                if err != nil {
                        fmt.Println("Error reading connection: ", err.Error())
                        os.Exit(1)
                }

                rl := extractRequestLine(string(b))
                if rl.verb == "GET" && strings.HasPrefix(rl.target, "/echo/") {

                        cl := len(rl.target) - len("/echo/")
                        body := rl.target[len("/echo/"):]
                        response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", cl, body)
                        conn.Write([]byte(response))

                } else if rl.verb == "GET" && rl.target == "/" {
                        conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

                } else {
                        conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
                }
                conn.Close()
        }
}

func extractRequestLine(request string) RequestLine {
        req := strings.Split(request, "\r\n")
        requestLine := req[0]
        rl := strings.Split(requestLine, " ")
        return RequestLine{rl[0], rl[1], rl[2]}
}
