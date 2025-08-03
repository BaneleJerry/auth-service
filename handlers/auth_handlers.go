package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/BaneleJerry/auth-service/dto"
	"github.com/BaneleJerry/auth-service/services"
	"github.com/BaneleJerry/auth-service/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: service}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userReq dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := h.authService.RegisterUser(userReq)
	if err != nil {
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userReq dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := h.authService.LoginUser(userReq)
	if err != nil {
		http.Error(w, "Failed to Login user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT token (use user.ID.String() if GenerateJWT expects string)
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert user to UserProfileResponse
	userProfile := dto.UserProfileResponseFromUser(user)

	// Create a response DTO
	resp := map[string]any{
		"user":  userProfile,
		"token": token,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Here you could do token blacklist or session invalidation if implemented
	// For JWT, usually you tell client to delete token

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User " + claims.Subject + " logged out successfully"))
}
