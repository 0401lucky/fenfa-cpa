package model

import (
	"cpa-distribution/common"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if common.SqlDSN != "" {
		DB, err = gorm.Open(postgres.Open(common.SqlDSN), gormConfig)
	} else {
		dbDir := "data"
		os.MkdirAll(dbDir, 0755)
		dbPath := filepath.Join(dbDir, "cpa.db")
		DB, err = gorm.Open(sqlite.Open(dbPath), gormConfig)
		if sqlDB, e := DB.DB(); e == nil {
			sqlDB.Exec("PRAGMA journal_mode=WAL")
			sqlDB.Exec("PRAGMA busy_timeout=5000")
		}
	}
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	if sqlDB, err := DB.DB(); err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	migrate()
}

func migrate() {
	err := DB.AutoMigrate(
		&User{},
		&Token{},
		&RequestLog{},
		&IPBan{},
		&SystemSetting{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
}
