package resolver

import (
	"bytes"
	"io"
	"net/http"

	"github.com/isa0-gh/easydoh/internal/config"
)

func Resolve(dnsmessage []byte) ([]byte, error) {
	var body []byte

	req, err := http.NewRequest("POST", config.Conf.Resolver, bytes.NewReader(dnsmessage))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/dns-message")
	req.Header.Set("Content-Type", "application/dns-message")

	resp, err := config.Conf.Client.Do(req)
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
