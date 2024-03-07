package configs

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectToDB(config *Config) (*gorm.DB, error) {
	var err error
	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		config.DBHost, config.DBUsername, config.DBUserPassword, config.DBName, config.DBPort)

	DB, err = gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to postgres")
		return nil, err
	}

	log.Println("Successfully connected to postgres")
	return DB, nil
}
