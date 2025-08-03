package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BaneleJerry/auth-service/dto"
	"github.com/BaneleJerry/auth-service/services"
	"github.com/BaneleJerry/auth-service/utils"
	"github.com/google/uuid"
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
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//REFRESH TOKEN
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,                    // prevents JS access (XSS protection)
		Secure:   false,                    // set to true in production (requires HTTPS)
		Path:     "/auth/refresh-token",    // restrict cookie to refresh route
		MaxAge:   30 * 24 * 3600,          // 30 days
		SameSite: http.SameSiteStrictMode, // or Lax depending on your needs
	})

	// Convert user to UserProfileResponse
	userProfile := dto.UserProfileResponseFromUser(user)

	// Create a response DTO
	resp := map[string]any{
		"user":    userProfile,
		"token":   token,
		"message": "Login successful",
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

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true, // same as above
		Path:     "/api/refresh-token",
		MaxAge:   -1, // immediately expires the cookie
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User " + claims.Subject + " logged out successfully"))
}

func (h *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get refresh token from cookie or request body

	// print cookie for debugging
	fmt.Println("Cookies:", r.Cookies())

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token cookie missing", http.StatusUnauthorized)
		return
	}
	refreshToken := cookie.Value

	// 2. Validate the refresh token
	token, claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil || !token.Valid {
		http.Error(w, "Invalid refresh token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// 3. Extract user ID from claims
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		http.Error(w, "Invalid subject in token", http.StatusUnauthorized)
		return
	}

	// 4. (Optional) Check if token is revoked/blacklisted in DB/Redis here

	// 5. Generate a new access token
	accessToken, err := utils.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// (Optional) generate new refresh token here if implementing refresh rotation
	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "refresh_token",
	// 	Value:    refreshToken,
	// 	HttpOnly: true,                    // prevents JS access (XSS protection)
	// 	Secure:   true,                    // set to true in production (requires HTTPS)
	// 	Path:     "/api/refresh-token",    // restrict cookie to refresh route
	// 	MaxAge:   30 * 24 * 3600,          // 30 days
	// 	SameSite: http.SameSiteStrictMode, // or Lax depending on your needs
	// })

	// 6. Return new tokens
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}
