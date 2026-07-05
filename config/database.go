package config

import (
	"log"
	"os"
	"spotsync-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a global variable that holds our database connection.
// We'll pass this into our repositories later.
var DB *gorm.DB

// ConnectDatabase reads the DATABASE_URL from environment variables
// and establishes a connection to PostgreSQL using GORM.
func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set in .env file")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("✅ Database connected successfully")

	// Auto-migrate: creates/updates tables based on our models
	err = db.AutoMigrate(&models.User{}, &models.ParkingZone{}, &models.Reservation{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	DB = db
}