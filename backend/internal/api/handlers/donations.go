package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"sustainwear/internal/api/middleware"
	"sustainwear/internal/config"
	"sustainwear/internal/domain/donation"
	"sustainwear/internal/domain/inventory"
	"sustainwear/internal/domain/user"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
)

type DonationHandler struct {
	donationService  *donation.Service
	inventoryService *inventory.Service
	userService      *user.Service
	config           *config.Config
}

func NewDonationHandler(donationService *donation.Service, inventoryService *inventory.Service, userService *user.Service, cfg *config.Config) *DonationHandler {
	return &DonationHandler{
		donationService:  donationService,
		inventoryService: inventoryService,
		userService:      userService,
		config:           cfg,
	}
}

// CREATE DONATION
func (h *DonationHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req donation.CreateRequest
	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("DONATIONS: [POST api/donations] - Failed to decode request body - %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	created, err := h.donationService.Create(userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "server error") {
			log.Printf("DONATIONS: [POST api/donations] - Failed to create donation - %v", err)
			http.Error(w, "Unable to create donation", http.StatusInternalServerError)
		} else {
			log.Printf("DONATIONS: [POST api/donations] - Bad request - %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsoniter.NewEncoder(w).Encode(created)
}

// GET DONATION BY ID
func (h *DonationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("DONATIONS: [GET api/donations/{id}] - Donation ID missing in request")
		http.Error(w, "Donation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("DONATIONS: [GET api/donations/{id}] - Invalid donation ID: %v", err)
		http.Error(w, "Invalid donation ID", http.StatusBadRequest)
		return
	}

	donation, err := h.donationService.GetByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("DONATIONS: [GET api/donations/{id}] - Donation not found with ID: %d", id)
			http.Error(w, "Donation not found", http.StatusNotFound)
		} else {
			log.Printf("DONATIONS: [GET api/donations/{id}] - Failed to get donation: %v", err)
			http.Error(w, "Unable to retrieve donation", http.StatusInternalServerError)
		}
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
		log.Printf("DONATIONS: [GET api/donations] - Failed to list donations: %v", err)
		http.Error(w, "Unable to list donations", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":      donations,
		"page":      page,
		"page_size": len(donations),
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
		log.Printf("DONATIONS: [GET api/donations/my] - Failed to list donations for user %d: %v", userID, err)
		http.Error(w, "Unable to get my donations", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":      donations,
		"page":      page,
		"page_size": len(donations),
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(response)
}

// UPDATE DONATION STATUS (CHARITY STAFF)
func (h *DonationHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("DONATIONS: [PUT api/donations/{id}/status] - Donation ID is required")
		http.Error(w, "Donation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("DONATIONS: [PUT api/donations/{id}/status] - Invalid donation ID: %v", err)
		http.Error(w, "Invalid donation ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}

	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("DONATIONS: [PUT api/donations/{id}/status] - Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.donationService.UpdateStatus(uint(id), req.Status, req.Notes)
	if err != nil {
		log.Printf("DONATIONS: [PUT api/donations/{id}/status] - Failed to update status for donation %d: %v", id, err)
		http.Error(w, "Unable to update donation status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(map[string]string{
		"message": "Status updated successfully",
	})
}

// APPROVE DONATION (CHARITY STAFF)
func (h *DonationHandler) Approve(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("DONATIONS: [POST api/donations/{id}/approve] - Donation ID is required")
		http.Error(w, "Donation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("DONATIONS: [POST api/donations/{id}/approve] - Invalid donation ID: %v", err)
		http.Error(w, "Invalid donation ID", http.StatusBadRequest)
		return
	}

	orgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID)
	if err != nil {
		log.Printf("DONATIONS: [POST api/donations/{id}/approve] - Unauthorized: %v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	donation, err := h.donationService.GetByID(uint(id))
	if err != nil {
		log.Printf("DONATIONS: [POST api/donations/{id}/approve] - Donation not found: %v", err)
		http.Error(w, "Donation not found", http.StatusNotFound)
		return
	}

	// PREVENT PANIC
	if donation.OrgID == nil {
		log.Printf("DONATIONS: [POST api/donations/%d/approve] - Donation has no organisation assigned", id)
		http.Error(w, "Donation has no organisation assigned", http.StatusBadRequest)
		return
	}

	if *donation.OrgID != orgID {
		log.Printf("DONATIONS: [POST api/donations/%d/approve] - Donation belongs to org %d, not %d", id, *donation.OrgID, orgID)
		http.Error(w, "Access denied: donation belongs to different organisation", http.StatusForbidden)
		return
	}

	err = h.donationService.UpdateStatus(uint(id), "approved", "Donation approved")
	if err != nil {
		log.Printf("DONATIONS: [POST api/donations/%d/approve] - Failed to approve: %v", id, err)
		http.Error(w, "Unable to approve donation", http.StatusInternalServerError)
		return
	}

	_, err = h.inventoryService.CreateFromDonation(
		donation.ID,
		donation.ItemName,
		donation.Category,
		donation.Condition,
		donation.Quantity,
		uint(orgID),
	)
	if err != nil {
		log.Printf("DONATIONS: [POST api/donations/%d/approve] - Failed to create inventory: %v", id, err)

		rollbackErr := h.donationService.UpdateStatus(uint(id), "pending", "Approval failed - inventory creation error")
		if rollbackErr != nil {
			log.Printf("DONATIONS: [POST api/donations/{id}/approve] - CRITICAL: Failed to rollback donation status: %v", rollbackErr)
			http.Error(w, "Critical error: donation approved but inventory failed. Contact support.", http.StatusInternalServerError)
			return
		}
		log.Printf("DONATIONS: [POST api/donations/%d/approve] - Rolled back donation %d to pending status due to inventory creation failure", id, id)
		http.Error(w, "Error creating inventory from donation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(map[string]string{
		"message": "Donation approved and added to inventory",
	})
}

// REJECT DONATION (CHARITY STAFF)
func (h *DonationHandler) Reject(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("DONATIONS: [POST api/donations/{id}/reject] - Donation ID is required")
		http.Error(w, "Donation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("DONATIONS: [POST api/donations/{id}/reject] - Invalid donation ID: %v", err)
		http.Error(w, "Invalid donation ID", http.StatusBadRequest)
		return
	}

	orgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID)
	if err != nil {
		log.Printf("DONATIONS: [POST api/donations/{id}/reject] - Unauthorized: %v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	donation, err := h.donationService.GetByID(uint(id))
	if err != nil {
		log.Printf("DONATIONS: [POST api/donations/{id}/reject] - Donation not found: %v", err)
		http.Error(w, "Donation not found", http.StatusNotFound)
		return
	}

	// PREVENT PANIC
	if donation.OrgID == nil {
		log.Printf("DONATIONS: [POST api/donations/%d/reject] - Donation has no organisation assigned", id)
		http.Error(w, "Donation has no organisation assigned", http.StatusBadRequest)
		return
	}

	if *donation.OrgID != orgID {
		log.Printf("DONATIONS: [POST api/donations/%d/reject] - Access denied: donation belongs to org %d, not %d", id, *donation.OrgID, orgID)
		http.Error(w, "Access denied: donation belongs to different organisation", http.StatusForbidden)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}

	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("DONATIONS: [POST api/donations/{id}/reject] - Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.donationService.UpdateStatus(uint(id), "rejected", req.Reason)
	if err != nil {
		log.Printf("DONATIONS: [POST api/donations/{id}/reject] - Failed to reject donation %d: %v", id, err)
		http.Error(w, "Unable to reject donation", http.StatusInternalServerError)
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
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("DONATIONS: [DELETE api/donations/{id}] - Donation ID is required")
		http.Error(w, "Donation ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("DONATIONS: [DELETE api/donations/{id}] - Invalid donation ID: %v", err)
		http.Error(w, "Invalid donation ID", http.StatusBadRequest)
		return
	}

	err = h.donationService.Delete(uint(id))
	if err != nil {
		log.Printf("DONATIONS: [DELETE api/donations/{id}] - Failed to delete donation %d: %v", id, err)
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
