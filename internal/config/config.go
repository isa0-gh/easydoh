package config

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Resolver    string            `toml:"resolver"`
	TTL         int               `toml:"ttl"`
	BindAddress string            `toml:"bind_address"`
	Hosts       map[string]string `toml:"hosts"`
	Client      *http.Client      `toml:"-"`
}

const DefaultConfigPath = "config.toml"

func Default() *Config {
	return &Config{
		Resolver:    "https://one.one.one.one/dns-query",
		TTL:         300,
		BindAddress: "0.0.0.0:53",
		Hosts: map[string]string{
			"*.home": "127.0.0.1",
		},
	}
}

func Load(path string) (*Config, error) {
	conf := Default()

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Info("Config file not found, using defaults", "path", path)
			if err := conf.Save(path); err != nil {
				slog.Error("Failed to save default config", "error", err)
			}
			return conf, nil
		}
		return nil, err
	}
	defer file.Close()

	if err := toml.NewDecoder(file).Decode(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Config) Save(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	encoder.SetIndentTables(true)
	return encoder.Encode(c)
}
