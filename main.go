package main

import (
	"fmt"
	"log"
	"net"

	"github.com/isa0-gh/easydoh/config"
	"github.com/isa0-gh/easydoh/providers"
	"github.com/isa0-gh/easydoh/resolver"
)

var bindAddr = config.Conf.BindAddress

func HandleConn(data []byte, addr *net.UDPAddr, conn *net.UDPConn) {

	resp, err := resolver.Resolve(data, providers.DnsProviders[config.Conf.Resolver])
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

	log.Printf("[+] easydoh started. Resolver: %s Listening on %s\n", config.Conf.Resolver, bindAddr)

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
