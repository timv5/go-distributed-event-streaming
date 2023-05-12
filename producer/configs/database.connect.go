package configs

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectToDB(config *Config) {
	var err error
	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		config.DBHost, config.DBUsername, config.DBUserPassword, config.DBName, config.DBPort)

	DB, err = gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to DB")
	} else {
		log.Println("Successfully connected to postgres")
	}
}
