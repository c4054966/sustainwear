package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"sustainwear/internal/api/middleware"
	"sustainwear/internal/config"
	"sustainwear/internal/domain/organisation"
	"sustainwear/internal/domain/user"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
)

type OrganisationHandler struct {
	organisationService *organisation.Service
	userService         *user.Service
	config              *config.Config
}

func NewOrganisationHandler(organisationService *organisation.Service, userService *user.Service, cfg *config.Config) *OrganisationHandler {
	return &OrganisationHandler{
		organisationService: organisationService,
		userService:         userService,
		config:              cfg,
	}
}

// CREATE ORGANISATION
func (h *OrganisationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req organisation.CreateOrgRequest
	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ORGANISATIONS: [POST api/organisations] - Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	created, err := h.organisationService.Create(&req)
	if err != nil {
		if strings.Contains(err.Error(), "server error") {
			log.Printf("ORGANISATIONS: [POST api/organisations] - Failed to create organisation: %v", err)
			http.Error(w, "Unable to create organisation", http.StatusInternalServerError)
		} else {
			log.Printf("ORGANISATIONS: [POST api/organisations] - Bad request: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsoniter.NewEncoder(w).Encode(created)
}

// GET ORGANISATION BY ID
func (h *OrganisationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("ORGANISATIONS: [GET api/organisations/{id}] - Missing organisation ID")
		http.Error(w, "Organisation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("ORGANISATIONS: [GET api/organisations/%s] - Invalid organisation ID value", idStr)
		http.Error(w, "Invalid organisation ID", http.StatusBadRequest)
		return
	}

	if role := middleware.GetUserRole(r); role == "charity_staff" {
		if allowedOrgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID); err != nil || allowedOrgID != uint(id) {
			log.Printf("ORGANISATIONS: [GET api/organisations/%d] - Unauthorized access attempt by charity_staff for org_id=%d", id, allowedOrgID)
			http.Error(w, "Access denied: can only access your own organisation", http.StatusForbidden)
			return
		}
	}

	org, err := h.organisationService.GetByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("ORGANISATIONS: [GET api/organisations/%d] - Organisation not found with ID: %d", id, id)
			http.Error(w, "Organisation not found", http.StatusNotFound)
		} else {
			log.Printf("ORGANISATIONS: [GET api/organisations/%d] - Failed to get organisation: %v", id, err)
			http.Error(w, "Unable to get organisation", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(org)
}

// GET ORGANISATION BY EMAIL
func (h *OrganisationHandler) GetByEmail(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]
	if email == "" {
		log.Printf("ORGANISATIONS: [GET api/organisations/email/{email}] - Email is missing in request")
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	org, err := h.organisationService.GetByEmail(email)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("ORGANISATIONS: [GET api/organisations/email/%s] - Organisation not found with email: %s", email, email)
			http.Error(w, "Organisation not found", http.StatusNotFound)
		} else {
			log.Printf("ORGANISATIONS: [GET api/organisations/email/%s] - Failed to get organisation with email: %s - %v", email, email, err)
			http.Error(w, "Unable to get organisation", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(org)
}

// LIST ORGANISATIONS
func (h *OrganisationHandler) List(w http.ResponseWriter, r *http.Request) {
	page := h.getIntParam(r, "page", 1)
	pageSize := h.getIntParam(r, "page_size", h.config.Pagination.DefaultPageSize)

	if pageSize > h.config.Pagination.MaxPageSize {
		pageSize = h.config.Pagination.MaxPageSize
	}

	offset := (page - 1) * pageSize

	filters := make(map[string]interface{})
	if orgType := r.URL.Query().Get("type"); orgType != "" {
		filters["type"] = orgType
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}
	if city := r.URL.Query().Get("city"); city != "" {
		filters["city"] = city
	}
	if county := r.URL.Query().Get("county"); county != "" {
		filters["county"] = county
	}

	filters["limit"] = pageSize
	filters["offset"] = offset

	orgs, err := h.organisationService.List(filters)
	if err != nil {
		log.Printf("ORGANISATIONS: [GET api/organisations] - Failed to list organisations: %v", err)
		http.Error(w, "Unable to list organisations", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":      orgs,
		"page":      page,
		"page_size": len(orgs),
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(response)
}

// UPDATE ORGANISATION
func (h *OrganisationHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("ORGANISATIONS: [PUT api/organisations/{id}] - Organisation ID is missing in request")
		http.Error(w, "Organisation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("ORGANISATIONS: [PUT api/organisations/%s] - Invalid organisation ID value", idStr)
		http.Error(w, "Invalid organisation ID", http.StatusBadRequest)
		return
	}

	if role := middleware.GetUserRole(r); role == "charity_staff" {
		if allowedOrgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID); err != nil || allowedOrgID != uint(id) {
			log.Printf("ORGANISATIONS: [PUT api/organisations/%d] - Unauthorized update attempt by charity_staff for org_id=%d", id, allowedOrgID)
			http.Error(w, "Access denied: can only access your own organisation", http.StatusForbidden)
			return
		}
	}

	var req organisation.UpdateOrgRequest
	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ORGANISATIONS: [PUT api/organisations/%d] - Invalid request body: %v", id, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updated, err := h.organisationService.Update(uint(id), &req)
	if err != nil {
		log.Printf("ORGANISATIONS: [PUT api/organisations/%d] - Failed to update organisation: %v", id, err)
		http.Error(w, "Unable to update organisation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(updated)
}

// DELETE ORGANISATION
func (h *OrganisationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("ORGANISATIONS: [DELETE api/organisations/{id}] - Organisation ID is missing in request")
		http.Error(w, "Organisation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("ORGANISATIONS: [DELETE api/organisations/%s] - Invalid organisation ID value", idStr)
		http.Error(w, "Invalid organisation ID", http.StatusBadRequest)
		return
	}

	if role := middleware.GetUserRole(r); role == "charity_staff" {
		if allowedOrgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID); err != nil || allowedOrgID != uint(id) {
			log.Printf("ORGANISATIONS: [DELETE api/organisations/%d] - Unauthorized delete attempt by charity_staff for org_id=%d", id, allowedOrgID)
			http.Error(w, "Access denied: can only access your own organisation", http.StatusForbidden)
			return
		}
	}

	err = h.organisationService.Delete(uint(id))
	if err != nil {
		log.Printf("ORGANISATIONS: [DELETE api/organisations/%d] - Failed to delete organisation: %v", id, err)
		http.Error(w, "Unable to delete organisation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET ORGANISATION STATS
func (h *OrganisationHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("ORGANISATIONS: [GET api/organisations/%s/stats] - Organisation ID is missing in request", idStr)
		http.Error(w, "Organisation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("ORGANISATIONS: [GET api/organisations/%s/stats] - Invalid organisation ID value", idStr)
		http.Error(w, "Invalid organisation ID", http.StatusBadRequest)
		return
	}

	if role := middleware.GetUserRole(r); role == "charity_staff" {
		if allowedOrgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID); err != nil || allowedOrgID != uint(id) {
			log.Printf("ORGANISATIONS: [GET api/organisations/%d/stats] - Unauthorized stats access attempt by charity_staff for org_id=%d", id, allowedOrgID)
			http.Error(w, "Access denied: can only access your own organisation", http.StatusForbidden)
			return
		}
	}

	stats, err := h.organisationService.GetStats(uint(id))
	if err != nil {
		log.Printf("ORGANISATIONS: [GET api/organisations/%d/stats] - Failed to get organisation stats: %v", id, err)
		http.Error(w, "Unable to get organisation stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(stats)
}

func (h *OrganisationHandler) getIntParam(r *http.Request, key string, defaultVal int) int {
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
