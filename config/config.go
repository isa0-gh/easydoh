package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Resolver string `json:"resolver"`
	TTL      int    `json:"ttl"`
}

var Conf Config

const configPath = "/etc/easydoh/config.json"

func init() {
	file, err := os.Open(configPath)
	if err != nil {
		Conf = Config{
			Resolver: "cloudflare",
			TTL:      300,
		}

		if err := saveConfig(); err != nil {
			panic("Failed to create config file: " + err.Error())
		}
		return
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&Conf); err != nil {
		Conf = Config{
			Resolver: "cloudflare",
			TTL:      300,
		}
		if err := saveConfig(); err != nil {
			panic("Failed to overwrite config file: " + err.Error())
		}
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
