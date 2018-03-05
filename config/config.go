package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

// Config represents database server and credentials
type Config struct {
	SMTPServer   string
	SMTPUsername string
	SMTPPassword string
	SMTPPort     int
}

// Read and parse the configuration file
func Read(configFile string) *Config {
	c := &Config{}
	if _, err := toml.DecodeFile(configFile, c); err != nil {
		log.Fatal(err)
	}
	return c
}
