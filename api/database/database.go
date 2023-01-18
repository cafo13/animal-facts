package database

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseHandler interface {
	GetDatabase() *Database
}

type Database struct {
	db *gorm.DB
}

func NewDatabaseHandler(dbHost string, dbPort string, dbName string, dbUser string, dbPassword string) (DatabaseHandler, error) {
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable TimeZone=UTC", dbHost, dbPort, dbName, dbUser, dbPassword)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("database connection error:", err)
	} else {
		log.Info("successfully connected to database")
	}

	db.AutoMigrate(&Fact{})

	return Database{db: db}, nil
}

func (db Database) GetDatabase() *Database {
	return &db
}
