package config

import "flag"

type Config struct {
	ServerAddress string
	BaseURL       string
}

func ParseFlags() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Base URL for shortened URLs")
	flag.Parse()
	return cfg
}
