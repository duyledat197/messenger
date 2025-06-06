package config

import (
	"embed"
	"log"

	"github.com/spf13/viper"
)

//go:embed config.yaml
var file embed.FS

// Config represents the overall configuration structure.

type Config struct {
	User struct {
		Postgres *Database `mapstructure:"postgres"`
		Endpoint *Endpoint `mapstructure:"endpoint"`
	} `mapstructure:"user"`

	Chat struct {
		OpenSearch *Database `mapstructure:"opensearch"`
		ScyllaDB   *Database `mapstructure:"scylla"`
		Courier    *Database `mapstructure:"courier"`
		Redis      *Database `mapstructure:"redis"`
		Endpoint   *Endpoint `mapstructure:"endpoint"`
	} `mapstructure:"chat"`

	Gateway struct {
		Endpoint *Endpoint `mapstructure:"endpoint"`
	}

	Security struct {
		SymetricKey string `mapstructure:"symetric_key"`
	} `mapstructure:"security"`

	Log struct {
		FileOutput string `mapstructure:"file_output"`
	} `mapstructure:"log"`
}

// LoadConfig loads the configuration from the specified file path and environment.
func LoadConfig() *Config {
	// Initialize an instance of the private config structure.
	var cfg Config
	f, err := file.Open("config.yaml")
	if err != nil {
		log.Fatalf("unable to open config file: %v", err)
	}

	viper.SetConfigType("yaml")
	// Read the configuration from the file.
	if err := viper.ReadConfig(f); err != nil {
		log.Fatalf("unable to read config file: %w", err)
	}

	// Unmarshal the configuration into the private config structure.
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to unmarshal config file: %w", err)
	}

	return &cfg
}
