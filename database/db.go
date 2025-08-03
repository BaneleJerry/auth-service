package database

import (
	"time"

	"github.com/BaneleJerry/auth-service/models"
	"github.com/go-faker/faker/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var err error
	sqliteDialector := sqlite.Open("test.db")
	if sqliteDialector == nil {
		panic("failed to open database")
	}

	DB, err = gorm.Open(sqliteDialector, &gorm.Config{})
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
