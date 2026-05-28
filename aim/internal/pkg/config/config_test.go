package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `
server:
  gateway:
    host: "0.0.0.0"
    port: 8080
    ws_path: "/ws"
  user:
    grpc_port: 50051
  message:
    grpc_port: 50052
  relation:
    grpc_port: 50053
  ai:
    grpc_port: 50054
    http_port: 8081
  file:
    grpc_port: 50055
    http_port: 8082
database:
  driver: "mysql"
  dsn: "root:pass@tcp(127.0.0.1:3306)/aim"
  max_open_conns: 50
  max_idle_conns: 5
redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0
auth:
  jwt_secret: "test-secret"
  token_expire: 3600
ai:
  default_provider: "openai"
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.Gateway.Port != 8080 {
		t.Errorf("GatewayPort = %d, want 8080", cfg.Server.Gateway.Port)
	}
	if cfg.Database.Driver != "mysql" {
		t.Errorf("DB Driver = %q, want mysql", cfg.Database.Driver)
	}
	if cfg.Redis.Addr != "127.0.0.1:6379" {
		t.Errorf("Redis Addr = %q", cfg.Redis.Addr)
	}
	if cfg.Auth.JWTSecret != "test-secret" {
		t.Errorf("JWTSecret = %q, want test-secret", cfg.Auth.JWTSecret)
	}
}

func TestDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	os.WriteFile(path, []byte("database:\n  driver: sqlite\n  dsn: \":memory:\"\n"), 0644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Database.MaxOpenConns != 100 {
		t.Errorf("default MaxOpenConns = %d, want 100", cfg.Database.MaxOpenConns)
	}
	if cfg.Auth.TokenExpire != 7200 {
		t.Errorf("default TokenExpire = %d, want 7200", cfg.Auth.TokenExpire)
	}
}

func TestLoad_InvalidPath(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Error("Load() should error on missing file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	os.WriteFile(path, []byte("invalid: [[[yaml"), 0644)

	_, err := Load(path)
	if err == nil {
		t.Error("Load() should error on invalid YAML")
	}
}
