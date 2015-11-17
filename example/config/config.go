package config

import "github.com/corneldamian/gowebservice"

var cfg *Config

func GetConfig() *Config {
	if cfg == nil {
		cfg = &Config{}
	}

	return cfg
}

type Config struct {
	WebServer gowebservice.WebServer
}

func (c *Config) WebServerConfig() *gowebservice.WebServer {
	return &c.WebServer
}
