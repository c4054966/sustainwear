package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"sustainwear/internal/api/middleware"
	"sustainwear/internal/config"
	"sustainwear/internal/domain/user"
	"sustainwear/pkg/validator"

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
	// THE != NIL EMPTY CHECKS ARE BECAUSE IF THEY PASS IN A VARIABLE WITH AN EMPTY STRING, IT WOULD OVERRIDE THE EXISTING VALUE STILL
	currentUserID := middleware.GetUserID(r)
	currentUserRole := middleware.GetUserRole(r)

	var req struct {
		UserID   *uint   `json:"user_id"`
		FullName *string `json:"full_name"`
		Email    *string `json:"email"`
		Role     *string `json:"role"`
		OrgID    *uint   `json:"org_id"`
		IsActive *bool   `json:"is_active"`
	}

	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("USERS: [PUT api/users/profile] - Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	targetUserID := currentUserID

	if req.UserID != nil {
		if currentUserRole != "admin" {
			log.Printf("USERS: [PUT api/users/profile] - Non-admin user %d attempted to update user %d", currentUserID, *req.UserID)
			http.Error(w, "Insufficient permissions to update other users", http.StatusForbidden)
			return
		}
		targetUserID = *req.UserID
	}

	existingUser, err := h.userService.GetByID(targetUserID)
	if err != nil {
		log.Printf("USERS: [PUT api/users/profile] - Failed to get profile for user ID %d: %v", targetUserID, err)
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Unable to update user profile", http.StatusInternalServerError)
		}
		return
	}

	if req.FullName != nil {
		if validator.IsEmpty(*req.FullName) {
			log.Printf("USERS: [PUT api/users/profile] - Full name string cannot be empty")
			http.Error(w, "Full name string cannot be empty", http.StatusBadRequest)
			return
		}
		existingUser.FullName = *req.FullName
	}

	if currentUserRole == "admin" {
		if req.Email != nil {
			if validator.IsEmpty(*req.Email) {
				log.Printf("USERS: [PUT api/users/profile] - Email string cannot be empty")
				http.Error(w, "Email string cannot be empty", http.StatusBadRequest)
				return
			}
			if !validator.IsValidEmail(*req.Email) {
				log.Printf("USERS: [PUT api/users/profile] - Invalid email format: %s", *req.Email)
				http.Error(w, "Invalid email format", http.StatusBadRequest)
				return
			}
			existingUser.Email = *req.Email
		}

		if req.Role != nil {
			if validator.IsEmpty(*req.Role) {
				log.Printf("USERS: [PUT api/users/profile] - Role string cannot be empty")
				http.Error(w, "Role string cannot be empty", http.StatusBadRequest)
				return
			}
			if !validator.IsValidRole(*req.Role) {
				log.Printf("USERS: [PUT api/users/profile] - Invalid role: %s", *req.Role)
				http.Error(w, "Invalid role. Must be: donor, charity_staff, or admin", http.StatusBadRequest)
				return
			}
			existingUser.Role = *req.Role
		}

		if req.OrgID != nil {
			existingUser.OrganisationID = req.OrgID
		}

		if req.IsActive != nil {
			existingUser.IsActive = *req.IsActive
		}
	} else {
		if req.Email != nil || req.Role != nil || req.OrgID != nil || req.IsActive != nil {
			log.Printf("USERS: [PUT api/users/profile] - User %d attempted to update admin-only fields", currentUserID)
			http.Error(w, "Insufficient permissions to update these fields", http.StatusForbidden)
			return
		}
	}

	err = h.userService.Update(existingUser)
	if err != nil {
		log.Printf("USERS: [PUT api/users/profile] - Failed to update user %d: %v", targetUserID, err)
		http.Error(w, "Unable to update user profile", http.StatusInternalServerError)
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
		"page_size": len(users),
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
