package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/BaneleJerry/auth-service/dto"
	"github.com/BaneleJerry/auth-service/services"
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
