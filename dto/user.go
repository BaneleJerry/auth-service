package dto

import (
	"time"

	"github.com/BaneleJerry/auth-service/models"
	"github.com/google/uuid"
)

// Requests
type RegisterUserRequest struct {
	Email       string     `json:"email" binding:"required,email"`
	Password    string     `json:"password" binding:"required,min=8"`
	ConfirmPass string     `json:"confirm_password" binding:"required,eqfield=Password"`
	FirstName   string     `json:"first_name" binding:"required"`
	LastName    string     `json:"last_name" binding:"required"`
	PhoneNumber string     `json:"phone_number" binding:"required,e164"`
	Address     *string    `json:"address,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserProfileRequest struct {
	FirstName   string     `json:"first_name,omitempty"`
	LastName    string     `json:"last_name,omitempty"`
	PhoneNumber string     `json:"phone_number,omitempty"`
	Address     *string    `json:"address,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// Responses
type UserProfileResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	IsVerified  bool      `json:"is_verified"`
}

// Outgoing from server → client (after login)
type AuthResponse struct {
	User  UserProfileResponse `json:"user"`
	Token string              `json:"token"`
}

func UserProfileResponseFromUser(user *models.User) UserProfileResponse {
	return UserProfileResponse{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		IsVerified:  user.IsVerified,
	}
}
