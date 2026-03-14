package main

import (
	"time"
	"log/slog"
	"net"
	"os"

	"github.com/isa0-gh/resolv/internal/cache"
	"github.com/isa0-gh/resolv/internal/config"
	"github.com/isa0-gh/resolv/internal/local"
	"github.com/isa0-gh/resolv/internal/resolver"
)

var bindAddr = config.Conf.BindAddress
var cdb = cache.New()
func HandleConn(data []byte, addr *net.UDPAddr, conn *net.UDPConn) {
	if localResp, ok := local.Match(data); ok {
		_, err := conn.WriteToUDP(localResp, addr)
		if err != nil {
			slog.Error("ERROR writing local resp", "error", err)
		}
		return
	}

	var resp []byte
	cached, ok, err := cdb.Get(data)
	if err == nil && ok {
		resp = cached
	} else {
		resp, err = resolver.Resolve(data)
		if err != nil {
			slog.Error("ERROR resolving dns message", "error", err)
			return
		}
		if err := cdb.Add(data, resp); err != nil {
			slog.Error("CACHE ERROR", "error", err)
		}

	}

	_, err = conn.WriteToUDP(resp, addr)
	if err != nil {
		slog.Error("ERROR writing to udp", "error", err)
		return
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	addr, err := net.ResolveUDPAddr("udp", bindAddr)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	slog.Info("resolv started.", "resolver", config.Conf.Resolver, "listen", bindAddr)

	buf := make([]byte, 4096)
	cdb.StartFlusher(time.Duration(config.Conf.TTL) * time.Second)
	for {
		size, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			slog.Error("ERROR reading from UDP", "error", err)
			continue
		}

		request := make([]byte, size)
		copy(request, buf[:size])
		go HandleConn(request, clientAddr, conn)

	}

}
