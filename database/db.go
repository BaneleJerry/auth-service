package database

import (
	"fmt"
	"os"
	"time"

	"github.com/BaneleJerry/auth-service/models"
	"github.com/go-faker/faker/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var err error

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_TIMEZONE"),
	)

	postgresDialector := postgres.Open(dsn)
	if postgresDialector == nil {
		panic("failed to open database")
	}

	DB, err = gorm.Open(postgresDialector, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		panic("failed to migrate database")
	}

	return DB
}

func SeedDB() {
	if DB == nil {
		panic("database not connected")
	}

	FakerUsers := make([]models.User, 0)
	for range 10 {
		Address := faker.GetRealAddress().Address
		// dob := faker.Date()

		FakerUsers = append(FakerUsers, models.User{
			Email:        faker.Email(),
			FirstName:    faker.FirstName(),
			LastName:     faker.LastName(),
			PasswordHash: "password",
			PhoneNumber:  faker.Phonenumber(),
			Address:      &Address,
			DateOfBirth:  &time.Time{},
			IsVerified:   false,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			DeletedAt:    gorm.DeletedAt{},
		})
	}
	DB.Create(&FakerUsers)
}
