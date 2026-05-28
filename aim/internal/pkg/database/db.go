package database

import (
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database configuration.
type Config struct {
	Driver       string `yaml:"driver"`
	DSN          string `yaml:"dsn"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

// DefaultConfig returns a MySQL-based config for development.
func DefaultConfig() Config {
	return Config{
		Driver:       "mysql",
		DSN:          "root:password@tcp(127.0.0.1:3306)/aim?charset=utf8mb4&parseTime=True&loc=Local",
		MaxOpenConns: 100,
		MaxIdleConns: 10,
	}
}

// New creates a new gorm.DB connection.
func New(cfg Config) (*gorm.DB, error) {
	dialector, err := dialector(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func dialector(driver, dsn string) (gorm.Dialector, error) {
	switch driver {
	case "mysql":
		return mysql.Open(dsn), nil
	case "sqlite":
		return sqlite.Open(dsn), nil
	default:
		return mysql.Open(dsn), nil
	}
}

// DB wraps gorm.DB for convenience.
type DB struct {
	*gorm.DB
}

// NewDB creates a wrapped DB instance.
func NewDB(cfg Config) (*DB, error) {
	gdb, err := New(cfg)
	if err != nil {
		return nil, err
	}
	return &DB{gdb}, nil
}

// TestDB returns an in-memory SQLite database for testing.
func TestDB() (*DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
