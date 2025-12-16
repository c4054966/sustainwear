package handlers

import (
	"log"
	"net/http"
	"strings"

	"sustainwear/internal/config"
	"sustainwear/internal/domain/user"
	"sustainwear/pkg/jwt"

	jsoniter "github.com/json-iterator/go"
)

type AuthHandler struct {
	userService *user.Service
	config      *config.Config
}

func NewAuthHandler(userService *user.Service, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		config:      cfg,
	}
}

// REGISTER NEW USER
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req user.RegisterRequest
	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("AUTH: [POST api/auth/register] - Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newUser, err := h.userService.Register(&req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			log.Printf("AUTH: [POST api/auth/register] - Validation error: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if strings.Contains(err.Error(), "server error") {
			log.Printf("AUTH: [POST api/auth/register] - Failed to register user: %v", err)
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
		} else {
			log.Printf("AUTH: [POST api/auth/register] - Bad Request: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// GENERATE JWT TOKEN
	token, err := jwt.GenerateToken(newUser.ID, newUser.Email, newUser.Role, h.config.Security.JWTSecret, h.config.Security.JWTExpiryHours)
	if err != nil {
		log.Printf("AUTH: [POST api/auth/register] - Failed to generate session token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := user.AuthResponse{
		Token:    token,
		UserID:   newUser.ID,
		Role:     newUser.Role,
		FullName: newUser.FullName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsoniter.NewEncoder(w).Encode(response)
}

// LOGIN USER
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req user.LoginRequest
	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("AUTH: [POST api/auth/login] - Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	authenticatedUser, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		log.Printf("AUTH: [POST api/auth/login] - Authentication failed: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// GENERATE JWT TOKEN
	token, err := jwt.GenerateToken(authenticatedUser.ID, authenticatedUser.Email, authenticatedUser.Role, h.config.Security.JWTSecret, h.config.Security.JWTExpiryHours)
	if err != nil {
		log.Printf("AUTH: [POST api/auth/login] - Failed to generate session token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := user.AuthResponse{
		Token:    token,
		UserID:   authenticatedUser.ID,
		Role:     authenticatedUser.Role,
		FullName: authenticatedUser.FullName,
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(response)
}

// LOGOUT USER
// TO ADD - TOKEN BLACKLIST TO INVALIDATE SESSION TOKENS
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	jsoniter.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}

// REFRESH TOKEN
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}

	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("AUTH: [POST api/auth/refresh] - Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// VALIDATE EXISTING TOKEN
	claims, err := jwt.ValidateToken(req.Token, h.config.Security.JWTSecret)
	if err != nil {
		log.Printf("AUTH: [POST api/auth/refresh] - Invalid or expired token: %v", err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// GENERATE NEW TOKEN
	newToken, err := jwt.GenerateToken(claims.UserID, claims.Email, claims.Role, h.config.Security.JWTSecret, h.config.Security.JWTExpiryHours)
	if err != nil {
		log.Printf("AUTH: [POST api/auth/refresh] - Failed to generate new token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(map[string]string{
		"token": newToken,
	})
}
