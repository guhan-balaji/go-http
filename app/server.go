package main

import (
        "fmt"
        "net"
        "os"
        "strings"
)

type Request struct {
        // request line
        verb, target, protocol string

        // headers
        host, userAgent, accept string

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
                defer conn.Close()

                request := make([]byte, 1024)
                _, err = conn.Read(request)
                if err != nil {
                        fmt.Println("Error reading connection: ", err.Error())
                        os.Exit(1)
                }

                req := deserializeHttpRequest(request)
                fmt.Println(req.verb)
                fmt.Println(req.target)
                fmt.Println(req.protocol)
                fmt.Println(req.host)
                fmt.Println(req.userAgent)
                fmt.Println(req.accept)
                sendResponse(req, conn)
                conn.Close()
        }
}

func deserializeHttpRequest(request []byte) *Request {
        reqParts := strings.Split(string(request), "\r\n")
        fmt.Println(reqParts)

        req := new(Request)

        reqLine := strings.Split(reqParts[0], " ")
        req.verb     = reqLine[0]
        req.target   = reqLine[1]
        req.protocol = reqLine[2]


        for _, rp := range reqParts[1:] {
                if strings.HasPrefix(rp, "Host") {
                        req.host = rp[len("Host: "):]
                } else if strings.HasPrefix(rp, "User-Agent") {
                        req.userAgent = rp[len("User-Agent: "):]
                } else if strings.HasPrefix(rp, "Accept") {
                        req.accept = rp[len("Accept: "):]
                }
        }
        return req
}

func sendResponse(req *Request, conn net.Conn) {

        if req.verb == "GET" && req.target == "/" {

                conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

        } else if req.verb == "GET" && strings.HasPrefix(req.target, "/echo/") {

                conn.Write([]byte(
                        fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
                        len(req.target) - len("/echo/"),
                        req.target[len("/echo/"):],
                )))

        } else if req.verb == "GET" && strings.HasPrefix(req.target, "/user-agent") {

                conn.Write([]byte(
                        fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
                        len(req.userAgent),
                        req.userAgent,
                )))

        } else {

                conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
        }
}
