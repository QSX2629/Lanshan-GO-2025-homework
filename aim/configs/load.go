// Package configs provides the shared configuration loader.
package configs

import (
	"log"

	pkgconfig "github.com/aim/aim/internal/pkg/config"
)

// Config is the top-level configuration.
type Config = pkgconfig.Config

// ServerConfig holds all server addresses.
type ServerConfig = pkgconfig.ServerConfig

// Load reads config from the default path.
func Load() *Config {
	cfg, err := pkgconfig.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("[config] load error: %v", err)
	}
	return cfg
}
