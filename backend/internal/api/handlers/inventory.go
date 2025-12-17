package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"sustainwear/internal/api/middleware"
	"sustainwear/internal/config"
	"sustainwear/internal/domain/inventory"
	"sustainwear/internal/domain/user"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
)

type InventoryHandler struct {
	inventoryService *inventory.Service
	userService      *user.Service
	config           *config.Config
}

func NewInventoryHandler(inventoryService *inventory.Service, userService *user.Service, cfg *config.Config) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
		userService:      userService,
		config:           cfg,
	}
}

// FOR MANUAL CREATION - SINCE INVENTORY IS HANDLED BY ALLOCATION OF DONATIONS
func (h *InventoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID)
	if err != nil {
		log.Printf("INVENTORY: [POST api/inventory] - Unauthorized access attempt - %v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var req inventory.CreateInventoryRequest
	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("INVENTORY: [POST api/inventory] - Failed to decode request body - %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	item, err := h.inventoryService.CreateManual(&req, uint(orgID))
	if err != nil {
		if strings.Contains(err.Error(), "server error") {
			log.Printf("INVENTORY: [POST api/inventory] - Failed to create inventory item - %v", err)
			http.Error(w, "Unable to create inventory item", http.StatusInternalServerError)
		} else {
			log.Printf("INVENTORY: [POST api/inventory] - Bad request - %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsoniter.NewEncoder(w).Encode(item)
}

// GET INVENTORY ITEM BY ID
func (h *InventoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("INVENTORY: [GET api/inventory/{id}] - Missing Inventory ID")
		http.Error(w, "Inventory ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("INVENTORY: [GET api/inventory/%s] - Invalid Inventory ID - %v", idStr, err)
		http.Error(w, "Invalid Inventory ID", http.StatusBadRequest)
		return
	}

	orgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID)
	if err != nil {
		log.Printf("INVENTORY: [GET api/inventory/%d] - Unauthorized access attempt - %v", id, err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	item, err := h.inventoryService.GetByID(uint(id), uint(orgID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("INVENTORY: [GET api/inventory/%d] - Inventory item not found with ID: %d", id, id)
			http.Error(w, "Inventory item not found", http.StatusNotFound)
		} else {
			log.Printf("INVENTORY: [GET api/inventory/%d] - Failed to retrieve inventory item - %v", id, err)
			http.Error(w, "Unable to find inventory item", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(item)
}

// LIST INVENTORY ITEMS
func (h *InventoryHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID)
	if err != nil {
		log.Printf("INVENTORY: [GET api/inventory] - Unauthorized access attempt - %v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	page := h.getIntParam(r, "page", 1)
	pageSize := h.getIntParam(r, "page_size", h.config.Pagination.DefaultPageSize)

	if pageSize > h.config.Pagination.MaxPageSize {
		pageSize = h.config.Pagination.MaxPageSize
	}

	offset := (page - 1) * pageSize

	filters := make(map[string]interface{})
	if category := r.URL.Query().Get("category"); category != "" {
		filters["category"] = category
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}
	if condition := r.URL.Query().Get("condition"); condition != "" {
		filters["condition"] = condition
	}

	filters["limit"] = pageSize
	filters["offset"] = offset

	items, err := h.inventoryService.List(uint(orgID), filters)
	if err != nil {
		log.Printf("INVENTORY: [GET api/inventory] - Failed to list inventory for org_id=%d - %v", orgID, err)
		http.Error(w, "Unable to load inventory", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":      items,
		"page":      page,
		"page_size": len(items),
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(response)
}

// UPDATE INVENTORY ITEM
func (h *InventoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("INVENTORY: [PUT api/inventory/{id}] - Missing Inventory ID")
		http.Error(w, "Inventory ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("INVENTORY: [PUT api/inventory/%s] - Invalid Inventory ID - %v", idStr, err)
		http.Error(w, "Invalid inventory ID", http.StatusBadRequest)
		return
	}

	orgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID)
	if err != nil {
		log.Printf("INVENTORY: [PUT api/inventory/%d] - Unauthorized access attempt - %v", id, err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var req inventory.UpdateInventoryRequest
	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("INVENTORY: [PUT api/inventory/%d] - Failed to decode request body - %v", id, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updated, err := h.inventoryService.Update(uint(id), uint(orgID), &req)
	if err != nil {
		log.Printf("INVENTORY: [PUT api/inventory/%d] - Failed to update inventory item for org_id=%d - %v", id, orgID, err)
		http.Error(w, "Unable to update inventory item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(updated)
}

// ALLOCATE INVENTORY
func (h *InventoryHandler) Allocate(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("INVENTORY: [POST api/inventory/{id}/allocate] - Missing Inventory ID")
		http.Error(w, "Inventory ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("INVENTORY: [POST api/inventory/%s/allocate] - Invalid Inventory ID - %v", idStr, err)
		http.Error(w, "Invalid inventory ID", http.StatusBadRequest)
		return
	}

	orgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID)
	if err != nil {
		log.Printf("INVENTORY: [POST api/inventory/%d/allocate] - Unauthorized access attempt - %v", id, err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("INVENTORY: [POST api/inventory/%d/allocate] - Failed to decode request body - %v", id, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.inventoryService.Allocate(uint(id), uint(orgID), req.Quantity)
	if err != nil {
		log.Printf("INVENTORY: [POST api/inventory/%d/allocate] - Failed to allocate inventory - %v", id, err)
		http.Error(w, "Unable to allocate inventory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(map[string]string{
		"message": "Inventory allocated successfully",
	})
}

// DISTRIBUTE INVENTORY
func (h *InventoryHandler) Distribute(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("INVENTORY: [POST api/inventory/{id}/distribute] - Missing Inventory ID")
		http.Error(w, "Inventory ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("INVENTORY: [POST api/inventory/%s/distribute] - Invalid Inventory ID - %v", idStr, err)
		http.Error(w, "Invalid inventory ID", http.StatusBadRequest)
		return
	}

	orgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID)
	if err != nil {
		log.Printf("INVENTORY: [POST api/inventory/%d/distribute] - Unauthorized access attempt - %v", id, err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("INVENTORY: [POST api/inventory/%d/distribute] - Failed to decode request body - %v", id, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.inventoryService.Distribute(uint(id), uint(orgID), req.Quantity)
	if err != nil {
		log.Printf("INVENTORY: [POST api/inventory/%d/distribute] - Failed to distribute inventory - %v", id, err)
		http.Error(w, "Unable to distribute inventory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(map[string]string{
		"message": "Inventory distributed successfully",
	})
}

// DEALLOCATE INVENTORY
func (h *InventoryHandler) Deallocate(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("INVENTORY: [POST api/inventory/{id}/deallocate] - Missing Inventory ID")
		http.Error(w, "Inventory ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("INVENTORY: [POST api/inventory/%s/deallocate] - Invalid Inventory ID - %v", idStr, err)
		http.Error(w, "Invalid inventory ID", http.StatusBadRequest)
		return
	}

	orgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID)
	if err != nil {
		log.Printf("INVENTORY: [POST api/inventory/%d/deallocate] - Unauthorized access attempt - %v", id, err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("INVENTORY: [POST api/inventory/%d/deallocate] - Failed to decode request body - %v", id, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.inventoryService.Deallocate(uint(id), uint(orgID), req.Quantity)
	if err != nil {
		log.Printf("INVENTORY: [POST api/inventory/%d/deallocate] - Failed to deallocate inventory - %v", id, err)
		http.Error(w, "Unable to deallocate inventory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(map[string]string{
		"message": "Inventory deallocated successfully",
	})
}

// DELETE INVENTORY ITEM
func (h *InventoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		log.Printf("INVENTORY: [DELETE api/inventory/{id}] - Missing Inventory ID")
		http.Error(w, "Inventory ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("INVENTORY: [DELETE api/inventory/%s] - Invalid Inventory ID - %v", idStr, err)
		http.Error(w, "Invalid inventory ID", http.StatusBadRequest)
		return
	}

	orgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID)
	if err != nil {
		log.Printf("INVENTORY: [DELETE api/inventory/%d] - Unauthorized access attempt - %v", id, err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = h.inventoryService.Delete(uint(id), uint(orgID))
	if err != nil {
		log.Printf("INVENTORY: [DELETE api/inventory/%d] - Failed to delete inventory item - %v", id, err)
		http.Error(w, "Unable to delete inventory item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET INVENTORY STATS
func (h *InventoryHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	orgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID)
	if err != nil {
		log.Printf("INVENTORY: [GET api/inventory/stats] - Unauthorized access attempt - %v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	stats, err := h.inventoryService.GetStats(uint(orgID))
	if err != nil {
		log.Printf("INVENTORY: [GET api/inventory/stats] - Failed to retrieve inventory stats - %v", err)
		http.Error(w, "Unable to retrieve inventory stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(stats)
}

func (h *InventoryHandler) getIntParam(r *http.Request, key string, defaultVal int) int {
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
