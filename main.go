package main

import (
	"fmt"
	"log"
	"net"

	"github.com/isa0-gh/cf-doh/resolver"
)

var bindAddr = "127.0.0.1:53"

func HandleConn(data []byte, addr *net.UDPAddr, conn *net.UDPConn) {

	resp, err := resolver.CloudflareDoH(data)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return
	}
	_, err = conn.WriteToUDP(resp, addr)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return
	}
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", bindAddr)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	log.Printf("[+] cf-doh started. Listening on %s\n", bindAddr)

	buf := make([]byte, 4096)

	for {
		size, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("ERROR: %s", err.Error())
			continue
		}

		request := make([]byte, size)
		copy(request, buf[:size])
		go HandleConn(request, clientAddr, conn)

	}

}
