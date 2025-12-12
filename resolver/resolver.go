package resolver

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

func CloudflareDoH(dnsmessage []byte) ([]byte, error) {
	var body []byte

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("POST", "https://1.1.1.1/dns-query", bytes.NewReader(dnsmessage))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/dns-message")
	req.Header.Set("Content-Type", "application/dns-message")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err

}
