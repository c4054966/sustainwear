package handlers

import (
	"net/http"
	"strconv"

	"sustainwear/internal/config"
	"sustainwear/internal/domain/inventory"

	jsoniter "github.com/json-iterator/go"
)

type InventoryHandler struct {
	inventoryService *inventory.Service
	config           *config.Config
}

func NewInventoryHandler(inventoryService *inventory.Service, cfg *config.Config) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
		config:           cfg,
	}
}

// GET INVENTORY ITEM BY ID
func (h *InventoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Inventory ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid inventory ID", http.StatusBadRequest)
		return
	}

	orgIDStr := r.URL.Query().Get("org_id")
	orgID, _ := strconv.ParseUint(orgIDStr, 10, 32)

	item, err := h.inventoryService.GetByID(uint(id), uint(orgID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(item)
}

// LIST INVENTORY ITEMS
func (h *InventoryHandler) List(w http.ResponseWriter, r *http.Request) {
	orgIDStr := r.URL.Query().Get("org_id")
	orgID, _ := strconv.ParseUint(orgIDStr, 10, 32)

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":      items,
		"page":      page,
		"page_size": pageSize,
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(response)
}

// UPDATE INVENTORY ITEM
func (h *InventoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Inventory ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid inventory ID", http.StatusBadRequest)
		return
	}

	orgIDStr := r.URL.Query().Get("org_id")
	orgID, _ := strconv.ParseUint(orgIDStr, 10, 32)

	var req inventory.UpdateInventoryRequest
	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updated, err := h.inventoryService.Update(uint(id), uint(orgID), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(updated)
}

// ALLOCATE INVENTORY
func (h *InventoryHandler) Allocate(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Inventory ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid inventory ID", http.StatusBadRequest)
		return
	}

	orgIDStr := r.URL.Query().Get("org_id")
	orgID, _ := strconv.ParseUint(orgIDStr, 10, 32)

	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.inventoryService.Allocate(uint(id), uint(orgID), req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(map[string]string{
		"message": "Inventory allocated successfully",
	})
}

// DISTRIBUTE INVENTORY
func (h *InventoryHandler) Distribute(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Inventory ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid inventory ID", http.StatusBadRequest)
		return
	}

	orgIDStr := r.URL.Query().Get("org_id")
	orgID, _ := strconv.ParseUint(orgIDStr, 10, 32)

	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.inventoryService.Distribute(uint(id), uint(orgID), req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(map[string]string{
		"message": "Inventory distributed successfully",
	})
}

// DEALLOCATE INVENTORY
func (h *InventoryHandler) Deallocate(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Inventory ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid inventory ID", http.StatusBadRequest)
		return
	}

	orgIDStr := r.URL.Query().Get("org_id")
	orgID, _ := strconv.ParseUint(orgIDStr, 10, 32)

	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.inventoryService.Deallocate(uint(id), uint(orgID), req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(map[string]string{
		"message": "Inventory deallocated successfully",
	})
}

// DELETE INVENTORY ITEM
func (h *InventoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Inventory ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid inventory ID", http.StatusBadRequest)
		return
	}

	orgIDStr := r.URL.Query().Get("org_id")
	orgID, _ := strconv.ParseUint(orgIDStr, 10, 32)

	err = h.inventoryService.Delete(uint(id), uint(orgID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET INVENTORY STATS
func (h *InventoryHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	orgIDStr := r.URL.Query().Get("org_id")
	if orgIDStr == "" {
		http.Error(w, "Organisation ID is required", http.StatusBadRequest)
		return
	}

	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid organisation ID", http.StatusBadRequest)
		return
	}

	stats, err := h.inventoryService.GetStats(uint(orgID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
