package local

import (
	"net"
	"strings"

	"github.com/miekg/dns"
)

type Matcher struct {
	hosts map[string]string
	ttl   int
}

func NewMatcher(hosts map[string]string, ttl int) *Matcher {
	return &Matcher{
		hosts: hosts,
		ttl:   ttl,
	}
}

// Match checks the DNS message against configured local hosts using exact or wildcard matching.
// It returns a boolean indicating if it intercepted the request, and a byte response if true.
func (m *Matcher) Match(reqBytes []byte) ([]byte, bool) {
	if len(m.hosts) == 0 {
		return nil, false
	}

	req := new(dns.Msg)
	if err := req.Unpack(reqBytes); err != nil {
		return nil, false
	}

	if len(req.Question) == 0 {
		return nil, false
	}

	q := req.Question[0]
	name := strings.TrimSuffix(q.Name, ".")

	if q.Qclass != dns.ClassINET {
		return nil, false
	}

	ipStr, found := m.matchHost(name)
	if !found {
		return nil, false
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, false // Invalid IP in config
	}

	resp := new(dns.Msg)
	resp.SetReply(req)

	isIPv4 := ip.To4() != nil

	if isIPv4 && q.Qtype == dns.TypeA {
		rr := new(dns.A)
		rr.Hdr = dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: uint32(m.ttl)}
		rr.A = ip.To4()
		resp.Answer = append(resp.Answer, rr)
	} else if (!isIPv4) && q.Qtype == dns.TypeAAAA {
		rr := new(dns.AAAA)
		rr.Hdr = dns.RR_Header{Name: q.Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: uint32(m.ttl)}
		rr.AAAA = ip.To16()
		resp.Answer = append(resp.Answer, rr)
	} else {
		// Found matching domain but wrong record type (e.g. asked for AAAA but configured with IPv4 A record)
		// We still matched it, so we return an empty successful response, rather than resolving upstream.
		resp.Rcode = dns.RcodeSuccess
	}

	respBytes, err := resp.Pack()
	if err != nil {
		return nil, false
	}

	return respBytes, true
}

func (m *Matcher) matchHost(domain string) (string, bool) {
	// Exact match
	if ip, ok := m.hosts[domain]; ok {
		return ip, true
	}

	// Wildcard match
	for pattern, ip := range m.hosts {
		if strings.HasPrefix(pattern, "*.") {
			suffix := strings.TrimPrefix(pattern, "*.")
			// Ensure it matches example.com for *.example.com OR sub.example.com
			if domain == suffix || strings.HasSuffix(domain, "."+suffix) {
				return ip, true
			}
		}
	}
	return "", false
}
