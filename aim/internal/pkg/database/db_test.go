package database

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Driver != "mysql" {
		t.Errorf("default driver = %q, want mysql", cfg.Driver)
	}
	if cfg.MaxOpenConns != 100 {
		t.Errorf("MaxOpenConns = %d, want 100", cfg.MaxOpenConns)
	}
}

func TestDialectorMySQL(t *testing.T) {
	d, err := dialector("mysql", "user:pass@tcp(127.0.0.1:3306)/db?charset=utf8mb4&parseTime=True")
	if err != nil {
		t.Fatalf("dialector() error = %v", err)
	}
	if d == nil {
		t.Fatal("dialector() returned nil")
	}
}

func TestDialectorSQLite(t *testing.T) {
	d, err := dialector("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("dialector() error = %v", err)
	}
	if d == nil {
		t.Fatal("dialector() returned nil")
	}
}

func TestDialectorUnknown(t *testing.T) {
	d, err := dialector("unknown", "dsn")
	if err != nil {
		t.Fatalf("dialector() error = %v", err)
	}
	if d == nil {
		t.Fatal("dialector() for unknown driver returned nil")
	}
}

func TestTestDB(t *testing.T) {
	db, err := TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	if db.DB == nil {
		t.Fatal("TestDB returned nil underlying DB")
	}

	type testModel struct {
		ID   uint   `gorm:"primaryKey"`
		Name string `gorm:"size:100"`
	}
	if err := db.AutoMigrate(&testModel{}); err != nil {
		t.Fatalf("AutoMigrate error = %v", err)
	}
	if err := db.Create(&testModel{Name: "test"}).Error; err != nil {
		t.Fatalf("Create error = %v", err)
	}

	var result testModel
	if err := db.First(&result, "name = ?", "test").Error; err != nil {
		t.Fatalf("First error = %v", err)
	}
	if result.Name != "test" {
		t.Errorf("Name = %q, want test", result.Name)
	}
}
