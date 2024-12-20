package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mskovv/tg-bot-subaru96/internal/models"
	"github.com/mskovv/tg-bot-subaru96/internal/models/migrations"
	"gorm.io/gorm/logger"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDb() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл.")
	}

	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Yekaterinburg",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("connected")
	DB.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migrations")
	err = DB.AutoMigrate(&models.User{}, &models.Appointment{}, &models.CarDictionary{})
	if err != nil {
		log.Fatal("Failed execute migrate. \n", err)
	}

	migrations.AddInitialUser(DB)
	migrations.CreateCarDictionary(DB)
	log.Println("end migrations")

	return DB
}
