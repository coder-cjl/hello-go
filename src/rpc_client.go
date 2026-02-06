package main

import (
	"fmt"
	"net/rpc"
)

type RpcClient struct{}

func (c RpcClient) Start() {
	fmt.Println("RPC Client Starting...")

	// 连接到RPC服务器
	conn, err := rpc.DialHTTP("tcp", "localhost:8000")
	if err != nil {
		fmt.Println("Dial error:", err)
		return
	}
	defer conn.Close()

	// 调用Area方法
	ret := 0
	err2 := conn.Call("RpcRect.Area", RpcParams{Width: 5, Height: 3}, &ret)
	if err2 != nil {
		fmt.Println("Call error:", err2)
	}
	fmt.Println("Area:", ret)

	// 调用Perimeter方法
	err3 := conn.Call("RpcRect.Perimeter", RpcParams{Width: 5, Height: 3}, &ret)
	if err3 != nil {
		fmt.Println("Call error:", err3)
	}
	fmt.Println("Perimeter:", ret)
}
