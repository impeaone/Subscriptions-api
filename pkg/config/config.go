package config

import "agrigation_api/pkg/tools"

type Config struct {
	Port      int    `yaml:"Port"`
	IPAddress string `yaml:"IPAddress"`
}

func ReadConfig() (*Config, error) {
	port := tools.GetEnvAsInt("SERVER_PORT", 11682)
	ip := tools.GetEnv("SERVER_IP", "127.0.0.1")
	config := &Config{
		Port:      port,
		IPAddress: ip,
	}
	return config, nil
}
