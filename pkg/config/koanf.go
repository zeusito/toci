package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Configurations struct {
	Server   ServerConfigurations   `koanf:"server"`
	Database DatabaseConfigurations `koanf:"database"`
	Hasher   HasherConfigurations   `koanf:"hasher"`
	Auth     AuthConfigurations     `koanf:"auth"`
	Email    EmailConfigurations    `koanf:"email"`
}

type ServerConfigurations struct {
	Port string `koanf:"port"`
}

type DatabaseConfigurations struct {
	Enabled    bool   `koanf:"enabled"`
	Host       string `koanf:"host"`
	Port       int    `koanf:"port"`
	DbName     string `koanf:"db-name"`
	Username   string `koanf:"user"`
	Password   string `koanf:"password"`
	PoolSize   int    `koanf:"pool-size"`
	LogQueries bool   `koanf:"log-queries"`
}

type HasherConfigurations struct {
	SHASecret string `koanf:"sha-secret"`
}

type AuthConfigurations struct {
	DevMode bool `koanf:"dev-mode"`
}

type EmailConfigurations struct {
	Enabled   bool   `koanf:"enabled"`
	DevMode   bool   `koanf:"dev-mode"`
	ApiKey    string `koanf:"api-key"`
	TestEmail string `koanf:"test-email"`
	FromEmail string `koanf:"from-email"`
}

// LoadConfigurations Loads configurations depending upon the environment
func LoadConfigurations(path string) (*Configurations, error) {
	k := koanf.New(".")
	err := k.Load(file.Provider(path), toml.Parser())
	if err != nil {
		return nil, err
	}

	// Searches for env variables and will transform them into koanf format
	// e.g. SERVER_PORT variable will be server.port: value
	err = k.Load(env.Provider("", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(s), "_", ".")
	}), nil)
	if err != nil {
		return nil, err
	}

	var configuration Configurations

	err = k.Unmarshal("", &configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}
