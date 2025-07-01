package config

import (
	"fmt"
	"log"
	"os"

	"github.com/GitNinja36/wello-backend/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(" Error loading .env file: %v", err)
	}

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal(" DB_URL is missing in .env")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(" Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.DoctorProfile{},
		&models.AdminProfile{},
		&models.Appointment{},
		&models.MedicalCheck{},
		&models.Order{},
		&models.Review{},
	)
	if err != nil {
		log.Fatalf(" AutoMigration failed: %v", err)
	}

	DB = db
	fmt.Println("Connected to DB & AutoMigrated successfully.")
}
