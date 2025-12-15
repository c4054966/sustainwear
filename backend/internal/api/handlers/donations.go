package handlers

import (
	"net/http"
	"strconv"

	"sustainwear/internal/api/middleware"
	"sustainwear/internal/config"
	"sustainwear/internal/domain/donation"

	jsoniter "github.com/json-iterator/go"
)

type DonationHandler struct {
	donationService *donation.Service
	config          *config.Config
}

func NewDonationHandler(donationService *donation.Service, cfg *config.Config) *DonationHandler {
	return &DonationHandler{
		donationService: donationService,
		config:          cfg,
	}
}

// CREATE DONATION
func (h *DonationHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req donation.CreateRequest
	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	created, err := h.donationService.Create(userID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsoniter.NewEncoder(w).Encode(created)
}

// GET DONATION BY ID
func (h *DonationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Donation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid donation ID", http.StatusBadRequest)
		return
	}

	donation, err := h.donationService.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(donation)
}

// LIST DONATIONS
func (h *DonationHandler) List(w http.ResponseWriter, r *http.Request) {
	page := h.getIntParam(r, "page", 1)
	pageSize := h.getIntParam(r, "page_size", h.config.Pagination.DefaultPageSize)

	if pageSize > h.config.Pagination.MaxPageSize {
		pageSize = h.config.Pagination.MaxPageSize
	}

	offset := (page - 1) * pageSize

	filters := make(map[string]interface{})
	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}
	if orgIDStr := r.URL.Query().Get("org_id"); orgIDStr != "" {
		orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
		if err == nil {
			filters["organization_id"] = uint(orgID)
		}
	}
	if donorIDStr := r.URL.Query().Get("donor_id"); donorIDStr != "" {
		donorID, err := strconv.ParseUint(donorIDStr, 10, 32)
		if err == nil {
			filters["donor_id"] = uint(donorID)
		}
	}

	filters["limit"] = pageSize
	filters["offset"] = offset

	donations, err := h.donationService.List(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":      donations,
		"page":      page,
		"page_size": pageSize,
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(response)
}

// LIST DONOR'S DONATIONS
func (h *DonationHandler) GetMyDonations(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	page := h.getIntParam(r, "page", 1)
	pageSize := h.getIntParam(r, "page_size", h.config.Pagination.DefaultPageSize)

	if pageSize > h.config.Pagination.MaxPageSize {
		pageSize = h.config.Pagination.MaxPageSize
	}

	offset := (page - 1) * pageSize

	filters := map[string]interface{}{
		"donor_id": userID,
		"limit":    pageSize,
		"offset":   offset,
	}

	donations, err := h.donationService.List(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":      donations,
		"page":      page,
		"page_size": pageSize,
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(response)
}

// UPDATE DONATION STATUS (CHARITY STAFF)
func (h *DonationHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Donation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid donation ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}

	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.donationService.UpdateStatus(uint(id), req.Status, req.Notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(map[string]string{
		"message": "Status updated successfully",
	})
}

// APPROVE DONATION (CHARITY STAFF)
func (h *DonationHandler) Approve(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Donation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid donation ID", http.StatusBadRequest)
		return
	}

	err = h.donationService.UpdateStatus(uint(id), "approved", "Donation approved")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(map[string]string{
		"message": "Donation approved successfully",
	})
}

// REJECT DONATION (CHARITY STAFF)
func (h *DonationHandler) Reject(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Donation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid donation ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}

	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.donationService.UpdateStatus(uint(id), "rejected", req.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(map[string]string{
		"message": "Donation rejected",
		"reason":  req.Reason,
	})
}

// DELETE DONATION
func (h *DonationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Donation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid donation ID", http.StatusBadRequest)
		return
	}

	err = h.donationService.Delete(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DonationHandler) getIntParam(r *http.Request, key string, defaultVal int) int {
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
