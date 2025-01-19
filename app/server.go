package main

import (
        "fmt"
        "net"
        "os"
        "strings"
        "strconv"
)

type Request struct {
        // request line
        verb, target, protocol string

        // headers
        host, userAgent, accept,
        contentType, contentLength string

        //body
        body string
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

                go func(c net.Conn) {
                        request := make([]byte, 1024)
                        _, err = c.Read(request)
                        if err != nil {
                                fmt.Println("Error reading connection: ", err.Error())
                                os.Exit(1)
                        }

                        req := deserializeHttpRequest(request)
                        sendResponse(req, c)
                        c.Close()
                }(conn)
        }
}

func deserializeHttpRequest(request []byte) *Request {
        var req Request
        reqParts := strings.Split(string(request), "\r\n")

        reqLine := strings.Split(reqParts[0], " ")
        req.verb     = reqLine[0]
        req.target   = reqLine[1]
        req.protocol = reqLine[2]


        for _, rp := range reqParts[1:] {
                switch {

                case strings.HasPrefix(rp, "Host"):
                        req.host = rp[len("Host: "):]
                case strings.HasPrefix(rp, "User-Agent"):
                        req.userAgent = rp[len("User-Agent: "):]
                case strings.HasPrefix(rp, "Accept"):
                        req.accept = rp[len("Accept: "):]
                case strings.HasPrefix(rp, "Content-Type: "):
                        req.contentType = rp[len("Content-Type: "):]
                case strings.HasPrefix(rp, "Content-Length: "):
                        req.contentLength = rp[len("Content-Length: "):]
                default:
                        continue

                }
        }

        if reqParts[len(reqParts) - 2] == "" {
                req.body = reqParts[len(reqParts) - 1]
        }

        return &req
}

func sendResponse(req *Request, conn net.Conn) {
        switch req.verb {
        case "GET":
                handleGetRequest(req, conn)
        case "POST":
                handlePostRequest(req, conn)
        default:
                conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))

        }
}

func handleGetRequest(req *Request, conn net.Conn) {
        if req.target == "/" {

                conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

        } else if strings.HasPrefix(req.target, "/echo/") {

                conn.Write([]byte(
                        fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
                        len(req.target) - len("/echo/"),
                        req.target[len("/echo/"):],
                )))

        } else if strings.HasPrefix(req.target, "/user-agent") {

                conn.Write([]byte(
                        fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
                        len(req.userAgent),
                        req.userAgent,
                )))

        } else if strings.HasPrefix(req.target, "/files") {
                fn := req.target[len("/files/"):]
                path := os.Args[2] + fn;
                content, err := os.ReadFile(path)

                if err != nil {
                        conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
                } else {

                        conn.Write([]byte(
                                fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s",
                                len(content),
                                content,
                        )))

                }

        } else {

                conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
        }
}

func handlePostRequest(req *Request, conn net.Conn) {
        if !strings.HasPrefix(req.target, "/files/") {
                conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
        }

        fn := req.target[len("/files/"):]
        path := os.Args[2] + fn
        l, err := strconv.Atoi(req.contentLength)
        if err != nil {
                conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
        }

        err = os.WriteFile(path, []byte(req.body[:l]), 0644)
        if err != nil {
                conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
        }
        conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
}
