// Package config provides configuration loading from YAML.
package config

import (
	"os"

	"github.com/aim/aim/internal/pkg/database"
	"gopkg.in/yaml.v3"
)

// Config is the top-level configuration.
type Config struct {
	Server   ServerConfig         `yaml:"server"`
	Database database.Config      `yaml:"database"`
	Redis    database.RedisConfig `yaml:"redis"`
	Auth     AuthConfig           `yaml:"auth"`
	AI       AIConfig             `yaml:"ai"`
}

// ServerConfig holds all server addresses.
type ServerConfig struct {
	Gateway  GatewayCfg    `yaml:"gateway"`
	User     GRPCCfg       `yaml:"user"`
	Message  GRPCCfg       `yaml:"message"`
	Relation GRPCCfg       `yaml:"relation"`
	AI       AIServerCfg   `yaml:"ai"`
	File     FileServerCfg `yaml:"file"`
}

// GatewayCfg holds the gateway server config.
type GatewayCfg struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	WSPath string `yaml:"ws_path"`
}

// GRPCCfg holds a gRPC server address.
type GRPCCfg struct {
	GRPCPort int `yaml:"grpc_port"`
}

// AIServerCfg holds AI service addresses.
type AIServerCfg struct {
	GRPCPort int `yaml:"grpc_port"`
	HTTPPort int `yaml:"http_port"`
}

// FileServerCfg holds file service addresses.
type FileServerCfg struct {
	GRPCPort int `yaml:"grpc_port"`
	HTTPPort int `yaml:"http_port"`
}

// AuthConfig holds authentication settings.
type AuthConfig struct {
	JWTSecret   string `yaml:"jwt_secret"`
	TokenExpire int    `yaml:"token_expire"`
}

// AIConfig holds AI-related settings.
type AIConfig struct {
	DefaultProvider      string          `yaml:"default_provider"`
	Providers            []AIProviderCfg `yaml:"providers"`
	MaxContextTokens     int             `yaml:"max_context_tokens"`
	SummaryPrompt        string          `yaml:"summary_prompt"`
	TodoExtractPrompt    string          `yaml:"todo_extract_prompt"`
	ReplyCandidatePrompt string          `yaml:"reply_candidate_prompt"`
}

// AIProviderCfg is a single AI provider entry.
type AIProviderCfg struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
}

// Load reads the config file and returns the parsed Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	setDefaults(cfg)
	return cfg, nil
}

func setDefaults(cfg *Config) {
	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 100
	}
	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 10
	}
	if cfg.Auth.TokenExpire == 0 {
		cfg.Auth.TokenExpire = 7200
	}
}
