package config

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/pelletier/go-toml/v2"

	"github.com/isa0-gh/resolv/internal/resolve-dns"
)

type Config struct {
	Resolver    string            `toml:"resolver"`
	TTL         int               `toml:"ttl"`
	BindAddress string            `toml:"bind_address"`
	Hosts       map[string]string `toml:"hosts"`
	Client      *http.Client
}

var Conf Config

const configPath = "/etc/resolv/config.toml"

func init() {
	// Try to open the config file
	file, err := os.Open(configPath)
	if err != nil {
		// File doesn't exist, create default config
		Conf = Config{
			Resolver:    "https://one.one.one.one/dns-query",
			TTL:         300,
			BindAddress: "0.0.0.0:53",
			Hosts: map[string]string{
				"*.home": "127.0.0.1",
			},
		}


		if err := saveConfig(); err != nil {
			panic("Failed to create config file: " + err.Error())
		}
		return
	}
	defer file.Close()

	if err := toml.NewDecoder(file).Decode(&Conf); err != nil {
		Conf = Config{
			Resolver:    "https://one.one.one.one/dns-query",
			TTL:         300,
			BindAddress: "0.0.0.0:53",
			Hosts: map[string]string{
				"*.home": "127.0.0.1",
			},
		}
		if err := saveConfig(); err != nil {
			panic("Failed to overwrite config file: " + err.Error())
		}
	}

	for {
		Conf.Client, err = resolvedns.ResolveServer(Conf.Resolver)
		if err != nil {
			slog.Error("Couldn't resolved server trying again...", "error", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
}

func saveConfig() error {
	if err := os.MkdirAll("/etc/resolv", 0755); err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	encoder.SetIndentTables(true)
	return encoder.Encode(Conf)
}
