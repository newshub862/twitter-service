package service

import (
	"fmt"
	"time"

	"newshub-twitter-service/dao"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func getDb() *gorm.DB {
	if db != nil {
		return db
	}

	if config.Driver == "sqlite3" {
		sqliteDB, err := gorm.Open(sqlite.Open(config.ConnectionString), &gorm.Config{})
		if err != nil {
			panic("open db error: " + err.Error())
		}

		db = sqliteDB
		return db
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DbPort,
	)

	pgdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("open db error: " + err.Error())
	}

	sqlDB, err := pgdb.DB()
	if err != nil {
		panic("open db error: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db = pgdb

	db.AutoMigrate(&dao.TwitterNews{})
	db.AutoMigrate(&dao.TwitterSource{})

	return db.Debug()
}
