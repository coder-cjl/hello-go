package main

import (
	"net/http"
)

type TestHttp struct{}

func (t TestHttp) Start() {
	http.HandleFunc("/go", goHandler)
	http.ListenAndServe("127.0.0.1:20001", nil)
}

// func (t TestHttp) Test() {
// 	resp, _ := http.Get("http://127.0.0.1:20001/go")
// 	defer resp.Body.Close()

// 	fmt.Println("Response status:", resp.Status)
// 	fmt.Println("Response headers:", resp.Header)

// 	buf, _ := make([]byte, 1024)
// 	for {
// 		n, err := resp.Body.Read(buf)
// 		if err != nil {
// 			break
// 		}
// 		fmt.Println("Response body:", string(buf[:n]))
// 	}
// }

func goHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, this is a response from Go HTTP server!"))
}
