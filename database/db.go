package database

import (
	"fmt"
	"goblog/config"
	"goblog/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	conf := config.AppConfig

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", conf.DBHost, conf.DBUser, conf.DBPassword, conf.DBName, conf.DBPort, conf.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database %v", err)
	}

	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Conment{}, &models.Like{}, &models.Tag{})

	DB = db
	log.Println("Database connected successfully!")
}
