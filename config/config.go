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
	PostgresDB *Database
	ScyllaDB   *Database
	OpenSearch *Database
	Courier    *Database

	UserService    *Endpoint
	ChatService    *Endpoint
	GatewayService *Endpoint

	SymetricKey   string
	FileLogOutPut string
}

// config is a private structure used for unmarshaling the configuration from Viper.
type config struct {
	DBHost          string `mapstructure:"DB_HOST"`
	DBPort          string `mapstructure:"DB_PORT"`
	DBUser          string `mapstructure:"DB_USER"`
	DBPassword      string `mapstructure:"DB_PASSWORD"`
	DBDatabase      string `mapstructure:"DB_NAME"`
	UserGRPCHost    string `mapstructure:"USER_GRPC_HOST"`
	UserGRPCPort    string `mapstructure:"USER_GRPC_PORT"`
	ChatGRPCHost    string `mapstructure:"CHAT_GRPC_HOST"`
	ChatGRPCPort    string `mapstructure:"CHAT_GRPC_PORT"`
	CouponGRPCHost  string `mapstructure:"COUPON_GRPC_HOST"`
	CouponGRPCPort  string `mapstructure:"COUPON_GRPC_PORT"`
	GatewayGRPCHost string `mapstructure:"GATEWAY_GRPC_HOST"`
	GatewayGRPCPort string `mapstructure:"GATEWAY_GRPC_PORT"`
	SymetricKey     string `mapstructure:"SYMETRIC_KEY"`
	FileLogOutPut   string `mapstructure:"FILE_LOG_OUTPUT"`
}

// LoadConfig loads the configuration from the specified file path and environment.
func LoadConfig() *Config {
	// Initialize an instance of the private config structure.
	var cfg config
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

	// Create and return the public Config structure based on the private config.
	return &Config{
		PostgresDB: &Database{
			Host:     cfg.DBHost,
			Port:     cfg.DBPort,
			User:     cfg.DBUser,
			Password: cfg.DBPassword,
			Database: cfg.DBDatabase,
		},
		UserService: &Endpoint{
			Host: cfg.UserGRPCHost,
			Port: cfg.UserGRPCPort,
		},
		ChatService: &Endpoint{
			Host: cfg.ChatGRPCHost,
			Port: cfg.ChatGRPCPort,
		},
		GatewayService: &Endpoint{
			Host: cfg.GatewayGRPCHost,
			Port: cfg.GatewayGRPCPort,
		},
		SymetricKey:   cfg.SymetricKey,
		FileLogOutPut: cfg.FileLogOutPut,
	}
}
