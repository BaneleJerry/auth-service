package main

import (
	"time"

	"github.com/BaneleJerry/auth-service/models"
	"github.com/go-faker/faker/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	gormDB, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = gormDB.AutoMigrate(&models.User{})
	if err != nil {
		panic("failed to migrate database")
	}

	FakerUsers := make([]models.User, 0)
	for range 10 {
		Address := faker.GetRealAddress().Address
		FakerUsers = append(FakerUsers, models.User{
			Email:        faker.Email(),
			PasswordHash: faker.Password(),
			FirstName:    faker.Name(),
			LastName:     faker.LastName(),
			PhoneNumber:  faker.Phonenumber(),
			Address:      &Address,
			DateOfBirth:  &time.Time{},
			IsVerified:   false,
		})
	}

	for range len(FakerUsers) {
		println("Generated User: " + FakerUsers[0].Email + " - " + FakerUsers[0].FirstName + " " + FakerUsers[0].LastName)
	}
	// gormDB.Create(&FakerUsers)

	// Print Users from the database
	var users []models.User
	gormDB.Find(&users)
	for _, user := range users {
		println(user.Email + " - " + user.FirstName + " " + user.LastName)
	}

}
