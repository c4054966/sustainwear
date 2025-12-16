package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"sustainwear/internal/api/middleware"
	"sustainwear/internal/config"
	"sustainwear/internal/domain/user"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
)

type UserHandler struct {
	userService *user.Service
	config      *config.Config
}

func NewUserHandler(userService *user.Service, cfg *config.Config) *UserHandler {
	return &UserHandler{
		userService: userService,
		config:      cfg,
	}
}

// GET USER PROFILE
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	profile, err := h.userService.GetByID(userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("USERS: [GET api/users/profile] - User profile not found for user ID %d", userID)
			http.Error(w, "User profile not found", http.StatusNotFound)
		} else {
			log.Printf("USERS: [GET api/users/profile] - Failed to get profile for user ID %d: %v", userID, err)
			http.Error(w, "Unable to get user profile", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(profile)
}

// UPDATE USER PROFILE
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	existingUser, err := h.userService.GetByID(userID)
	if err != nil {
		log.Printf("USERS: [PUT api/users/profile] - Failed to update profile for user ID %d: %v", userID, err)
		http.Error(w, "Unable to update user profile", http.StatusInternalServerError)
		return
	}

	var req struct {
		FullName *string `json:"full_name"`
	}

	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FullName != nil {
		existingUser.FullName = *req.FullName
	}

	err = h.userService.Update(existingUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(existingUser)
}

// GET USER BY ID (ADMIN/ORG STAFF)
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("USERS: [GET api/users/{id}] - User ID is missing in request")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("USERS: [GET api/users/{id}] - Invalid user ID: %s", idStr)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("USERS: [GET api/users/{id}] - User profile not found for user ID %d", id)
			http.Error(w, "User profile not found", http.StatusNotFound)
		} else {
			log.Printf("USERS: [GET api/users/{id}] - Failed to get profile for user ID %d: %v", id, err)
			http.Error(w, "Unable to get user profile", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(user)
}

// LIST ALL USERS (ADMIN/ORG STAFF)
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	page := h.getIntParam(r, "page", 1)
	pageSize := h.getIntParam(r, "page_size", h.config.Pagination.DefaultPageSize)

	if pageSize > h.config.Pagination.MaxPageSize {
		pageSize = h.config.Pagination.MaxPageSize
	}

	offset := (page - 1) * pageSize

	users, err := h.userService.ListPaginated(pageSize, offset)
	if err != nil {
		log.Printf("USERS: [GET api/users] - Failed to list users: %v", err)
		http.Error(w, "Unable to list users", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":      users,
		"page":      page,
		"page_size": pageSize,
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(response)
}

// DELETE USER (ADMIN ONLY)
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("USERS: [DELETE api/users/{id}] - User ID is missing in request")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("USERS: [DELETE api/users/{id}] - Invalid user ID: %s", idStr)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.userService.Delete(uint(id))
	if err != nil {
		log.Printf("USERS: [DELETE api/users/{id}] - Failed to delete user with ID %d: %v", id, err)
		http.Error(w, "Unable to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) getIntParam(r *http.Request, key string, defaultVal int) int {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil || val < 1 {
		return defaultVal
	}
	return val
}
