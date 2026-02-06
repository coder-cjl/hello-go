package main

import (
	"fmt"
	"net/http"
	"net/rpc"
)

type RpcParams struct {
	Width, Height int
}

type RpcRect struct{}

func (r RpcRect) Area(params RpcParams, reply *int) error {
	*reply = params.Width * params.Height
	return nil
}

func (r RpcRect) Perimeter(params RpcParams, reply *int) error {
	*reply = 2 * (params.Width + params.Height)
	return nil
}

type RpcService struct{}

func (p RpcService) Start() {

	fmt.Println("RPC Service Starting...")

	rect := new(RpcRect)
	rpc.Register(rect)
	rpc.HandleHTTP()

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("ListenAndServe error:", err)
	}
}
