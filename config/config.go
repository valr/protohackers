package config

import (
	"flag"
)

type Config struct {
	ServerHost string
	ServerPort string
}

const (
	defaultServerHost = "localhost"
	defaultServerPort = "12345"
)

func (cfg *Config) ParseFlags() {
	flag.StringVar(&cfg.ServerHost, "host", defaultServerHost, "Listen on the specified `host`")
	flag.StringVar(&cfg.ServerPort, "port", defaultServerPort, "Listen on the specified `port`")
	flag.Parse()
}
