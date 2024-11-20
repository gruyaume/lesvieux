package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type TLSYaml struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

type ConfigYAML struct {
	DBPath string  `yaml:"db_path"`
	Port   int     `yaml:"port"`
	TLS    TLSYaml `yaml:"tls"`
}

type TLS struct {
	Cert []byte
	Key  []byte
}

type Config struct {
	DBPath string
	Port   int
	TLS    TLS
}

func Validate(filePath string) (Config, error) {
	config := Config{}
	configYaml, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("cannot read config file: %w", err)
	}
	c := ConfigYAML{}
	if err := yaml.Unmarshal(configYaml, &c); err != nil {
		return Config{}, fmt.Errorf("cannot unmarshal config file: %w", err)
	}
	if c.TLS.Cert == "" {
		return Config{}, fmt.Errorf("tls.cert is empty")
	}
	cert, err := os.ReadFile(c.TLS.Cert)
	if err != nil {
		return Config{}, fmt.Errorf("cannot read cert file: %w", err)
	}
	if c.TLS.Key == "" {
		return Config{}, fmt.Errorf("tls.key is empty")
	}
	key, err := os.ReadFile(c.TLS.Key)
	if err != nil {
		return Config{}, fmt.Errorf("cannot read key file: %w", err)
	}
	if c.DBPath == "" {
		return Config{}, errors.New("`db_path` is empty")
	}
	dbfile, err := os.OpenFile(c.DBPath, os.O_CREATE|os.O_RDONLY, 0o644)
	if err != nil {
		return Config{}, fmt.Errorf("cannot open db file: %w", err)
	}
	err = dbfile.Close()
	if err != nil {
		return Config{}, fmt.Errorf("cannot close db file: %w", err)
	}
	if c.Port == 0 {
		return Config{}, errors.New("port is empty")
	}
	config.Port = c.Port
	config.TLS.Cert = cert
	config.TLS.Key = key
	config.DBPath = c.DBPath
	return config, nil
}
