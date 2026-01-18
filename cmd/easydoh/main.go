package main

import (
	"time"
	"fmt"
	"log"
	"net"

	"github.com/isa0-gh/easydoh/internal/cache"
	"github.com/isa0-gh/easydoh/internal/config"
	"github.com/isa0-gh/easydoh/internal/resolver"
)

var bindAddr = config.Conf.BindAddress
var cdb = cache.New()
func HandleConn(data []byte, addr *net.UDPAddr, conn *net.UDPConn) {
	var resp []byte
	cached, ok, err := cdb.Get(data)
	if err == nil && ok {
		resp = cached
	} else {
		resp, err = resolver.Resolve(data)
		if err != nil {
			fmt.Printf("ERROR: %s", err.Error())
			return
		}
		if err := cdb.Add(data, resp); err != nil {
			fmt.Println("[CACHE ERROR]:", err)
		}

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

	log.Printf("[+] easydoh started.\n[+] Resolver: %s Listening on %s\n", config.Conf.Resolver, bindAddr)

	buf := make([]byte, 4096)
	cdb.StartFlusher(time.Duration(config.Conf.TTL) * time.Second)
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
