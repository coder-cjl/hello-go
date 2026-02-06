package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type TcpScan struct{}

func isOpen(host string, port int, timeout time.Duration) bool {
	time.Sleep(time.Millisecond * 1)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err == nil {
		conn.Close()
		return true
	}
	return false
}

func (t TcpScan) Start() {
	ports := []int{}

	wg := &sync.WaitGroup{}
	timeout := time.Millisecond * 200

	for port := 1; port < 100; port++ {
		wg.Add(1)
		go func(p int) {
			opened := isOpen("localhost", p, timeout)
			if opened {
				ports = append(ports, p)
			}
			wg.Done()
		}(port)

	}
	wg.Wait()
	Log.Infof("opened ports: %v\n", ports)
}
