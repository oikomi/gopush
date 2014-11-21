package main

import (
		"fmt"
		"net"
       )

func main() {
	addrs, err := net.InterfaceAddrs()
	fmt.Println(addrs)
		if err != nil {
			panic(err)
		}
	for _, addr := range addrs {
		fmt.Println(addr.String())
	}
}
