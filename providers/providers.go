package providers

var DnsProviders map[string]string = map[string]string{
	"cloudflare": "https://1.1.1.1/dns-query",
	"google":     "https://8.8.8.8/dns-query",
	"quad9":      "https://9.9.9.9/dns-query",
	"cisco":      "https://208.67.222.222/dns-query",
	"adguard":    "https://94.140.14.14/dns-query",
}
