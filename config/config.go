package config

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/isa0-gh/easydoh/dns"
)

type Config struct {
	Resolver    string `json:"resolver"`
	TTL         int    `json:"ttl"`
	BindAddress string `json:"bind_address"`
	Client      *http.Client
}

var Conf Config

const configPath = "/etc/easydoh/config.json"

func init() {
	// Try to open the config file
	file, err := os.Open(configPath)
	if err != nil {
		// File doesn't exist, create default config
		Conf = Config{
			Resolver:    "https://one.one.one.one/dns-query",
			TTL:         300,
			BindAddress: "127.0.0.1:53",
		}

		if err := saveConfig(); err != nil {
			panic("Failed to create config file: " + err.Error())
		}
		return
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&Conf); err != nil {
		Conf = Config{
			Resolver:    "https://one.one.one.one/dns-query",
			TTL:         300,
			BindAddress: "127.0.0.1:53",
		}
		if err := saveConfig(); err != nil {
			panic("Failed to overwrite config file: " + err.Error())
		}
	}

	Conf.Client, err = dns.ResolveServer(Conf.Resolver)
	if err != nil {
		panic(err)
	}
}

func saveConfig() error {
	if err := os.MkdirAll("/etc/easydoh", 0755); err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(Conf)
}
