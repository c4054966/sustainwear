package handlers

import (
	"net/http"
	"strconv"

	"sustainwear/internal/config"
	"sustainwear/internal/domain/organisation"

	jsoniter "github.com/json-iterator/go"
)

type OrganisationHandler struct {
	organisationService *organisation.Service
	config              *config.Config
}

func NewOrganisationHandler(organisationService *organisation.Service, cfg *config.Config) *OrganisationHandler {
	return &OrganisationHandler{
		organisationService: organisationService,
		config:              cfg,
	}
}

// CREATE ORGANISATION
func (h *OrganisationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req organisation.CreateOrgRequest
	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	created, err := h.organisationService.Create(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsoniter.NewEncoder(w).Encode(created)
}

// GET ORGANISATION BY ID
func (h *OrganisationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Organisation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid organisation ID", http.StatusBadRequest)
		return
	}

	org, err := h.organisationService.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(org)
}

// GET ORGANISATION BY EMAIL
func (h *OrganisationHandler) GetByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	org, err := h.organisationService.GetByEmail(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":      orgs,
		"page":      page,
		"page_size": pageSize,
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(response)
}

// UPDATE ORGANISATION
func (h *OrganisationHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Organisation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid organisation ID", http.StatusBadRequest)
		return
	}

	var req organisation.UpdateOrgRequest
	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updated, err := h.organisationService.Update(uint(id), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(updated)
}

// DELETE ORGANISATION
func (h *OrganisationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Organisation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid organisation ID", http.StatusBadRequest)
		return
	}

	err = h.organisationService.Delete(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET ORGANISATION STATS
func (h *OrganisationHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Organisation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid organisation ID", http.StatusBadRequest)
		return
	}

	stats, err := h.organisationService.GetStats(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
