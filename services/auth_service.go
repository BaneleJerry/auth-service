package services

import (
	"fmt"
	"time"

	"github.com/BaneleJerry/auth-service/dto"
	"github.com/BaneleJerry/auth-service/models"
	"github.com/BaneleJerry/auth-service/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) RegisterUser(req dto.RegisterUserRequest) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PhoneNumber:  req.PhoneNumber,
		Address:      req.Address,
		DateOfBirth:  req.DateOfBirth,
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) LoginUser(req dto.LoginRequest) (*models.User, error) {
	user, err := s.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err 
	}

	if checked := utils.CheckPasswordHash(req.Password, user.PasswordHash); checked == false {
		// Password does not match
		return nil, fmt.Errorf("invalid email or password")
	}

	return user, nil
}

func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// return nil and a custom error your handler can interpret
			return nil, fmt.Errorf("user with email '%s' not found", email)
		}
		// log or wrap DB-related errors
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}
	return &user, nil
}
