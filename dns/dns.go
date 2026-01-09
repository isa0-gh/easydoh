package dns

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	NetURL "net/url"
	"time"
)

func GetIP(hostname string) (string, error) {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 5 * time.Second,
			}
			return d.DialContext(ctx, "udp", "1.1.1.1:53")
		},
	}

	ctx := context.Background()

	ips, err := resolver.LookupIP(ctx, "ip", hostname)
	if err != nil || len(ips) < 1 {
		return "", err
	}

	return ips[0].String(), nil

}

func ResolveServer(rawUrl string) (*http.Client, error) {
	// Resolves DNS server and using same ip always skips dns resolving
	url, err := NetURL.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	domain := url.Hostname()
	ip, err := GetIP(domain)
	if err != nil {
		return nil, err
	}

	dialer := &net.Dialer{Timeout: 10 * time.Second}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, _ := net.SplitHostPort(addr)
			if host == domain {
				addr = net.JoinHostPort(ip, port)
			}
			return dialer.DialContext(ctx, network, addr)
		},
		ForceAttemptHTTP2: true,
		TLSClientConfig:   &tls.Config{ServerName: domain},
	}

	return &http.Client{Transport: transport}, nil
}
