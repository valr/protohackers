package config

import (
	"flag"
)

type Config struct {
	ServerHost string
	ServerPort string
	DebugMode  bool
}

const (
	defaultServerHost = "localhost"
	defaultServerPort = "12345"
	defaultDebugMode  = false
)

func New() *Config {
	return &Config{defaultServerHost, defaultServerPort, defaultDebugMode}
}

func (cfg *Config) ParseFlags() {
	flag.StringVar(&cfg.ServerHost, "host", defaultServerHost, "Listen on the specified `host`")
	flag.StringVar(&cfg.ServerPort, "port", defaultServerPort, "Listen on the specified `port`")
	flag.BoolVar(&cfg.DebugMode, "debug", defaultDebugMode, "Display debugging output")
	flag.Parse()
}
