package database

import "testing"

func TestDefaultRedisConfig(t *testing.T) {
	cfg := DefaultRedisConfig()
	if cfg.Addr != "127.0.0.1:6379" {
		t.Errorf("Addr = %q, want 127.0.0.1:6379", cfg.Addr)
	}
	if cfg.DB != 0 {
		t.Errorf("DB = %d, want 0", cfg.DB)
	}
}
