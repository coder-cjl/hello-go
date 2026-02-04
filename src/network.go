package main

// import (
// 	"bufio"
// 	"fmt"
// 	"net"
// 	"strings"
// )

// func process(conn net.Conn) {
// 	defer conn.Close()

// 	reader := bufio.NewReader(conn)

// 	for {
// 		// 读取完整的 HTTP 请求（直到遇到空行）
// 		var request strings.Builder
// 		for {
// 			line, err := reader.ReadString('\n')
// 			if err != nil {
// 				fmt.Println("Read error:", err)
// 				return
// 			}
// 			request.WriteString(line)
// 			// HTTP 头部以 \r\n\r\n 结束，即空行
// 			if line == "\r\n" {
// 				break
// 			}
// 		}

// 		recvStr := request.String()
// 		fmt.Println("Received from client:\n", recvStr)

// 		// 返回 HTTP 响应
// 		body := "Hello from Go TCP Server!"
// 		response := "HTTP/1.1 200 OK\r\n" +
// 			"Content-Type: text/plain; charset=utf-8\r\n" +
// 			"Connection: close\r\n" +
// 			"Content-Length: " + fmt.Sprintf("%d", len(body)) + "\r\n" +
// 			"\r\n" +
// 			body
// 		conn.Write([]byte(response))
// 	}
// }

// func main() {

// 	fmt.Println("Starting server on")

// 	listener, err := net.Listen("tcp", "127.0.0.1:20000")
// 	if err != nil {
// 		fmt.Println("Listen error:", err)
// 		return
// 	}

// 	for {
// 		/// 建立连接
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			fmt.Println("Accept error:", err)
// 			continue
// 		}

// 		go process(conn)
// 	}
// }
