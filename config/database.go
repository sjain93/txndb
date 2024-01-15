package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := buildDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	DB = db
}

func buildDSN() string {
	// Options for SSL Mode are "disabled"/"enabled"
	return fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=%s",
		/*
			These variables can be found in a local .env file
			to run locally, init a postgres server and create a new
			database with a user that this project can use
		*/
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("PORT"),
		os.Getenv("SSL_mode"),
	)
}
