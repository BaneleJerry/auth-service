package models

import "time"

type User struct {
	ID          uint       `gorm:"primaryKey"`
	Email       string     `gorm:"unique;not null"`
	PasswordHash string     `gorm:"not null"`
	FirstName   string     `gorm:"not null"`
	LastName    string     `gorm:"not null"`
	PhoneNumber string     `gorm:"unique"`
	Address     *string    `gorm:"type:text"`
	DateOfBirth *time.Time `gorm:"type:date"`
	IsVerified  bool       `gorm:"not null;default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
