package handlers

import (
	"net/http"
	"strconv"
	"time"

	"sustainwear/internal/api/middleware"
	"sustainwear/internal/config"
	"sustainwear/internal/domain/analytics"

	jsoniter "github.com/json-iterator/go"
)

type AnalyticsHandler struct {
	analyticsService *analytics.Service
	config           *config.Config
}

func NewAnalyticsHandler(analyticsService *analytics.Service, cfg *config.Config) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		config:           cfg,
	}
}

// GET DONATION TRENDS
func (h *AnalyticsHandler) GetDonationTrends(w http.ResponseWriter, r *http.Request) {
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

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "daily"
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	startDate, err := time.Parse("02-01-2006", startDateStr)
	if err != nil {
		http.Error(w, "Invalid start_date format (use DD-MM-YYYY)", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("02-01-2006", endDateStr)
	if err != nil {
		http.Error(w, "Invalid end_date format (use DD-MM-YYYY)", http.StatusBadRequest)
		return
	}

	trends, err := h.analyticsService.GetDonationTrends(uint(orgID), period, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(trends)
}

// GET CATEGORY BREAKDOWN
func (h *AnalyticsHandler) GetCategoryBreakdown(w http.ResponseWriter, r *http.Request) {
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

	breakdown, err := h.analyticsService.GetCategoryBreakdown(uint(orgID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(breakdown)
}

// GET SUSTAINABILITY METRICS
func (h *AnalyticsHandler) GetSustainabilityMetrics(w http.ResponseWriter, r *http.Request) {
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

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "all_time"
	}

	metrics, err := h.analyticsService.GetSustainabilityMetrics(uint(orgID), period)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(metrics)
}

// GET DONOR IMPACT
func (h *AnalyticsHandler) GetDonorImpact(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	impact, err := h.analyticsService.GetDonorImpact(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(impact)
}

// GET ORG PERFORMANCE
func (h *AnalyticsHandler) GetOrgPerformance(w http.ResponseWriter, r *http.Request) {
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

	performance, err := h.analyticsService.GetOrgPerformance(uint(orgID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(performance)
}

// GET SYSTEM OVERVIEW (ADMIN ONLY)
func (h *AnalyticsHandler) GetSystemOverview(w http.ResponseWriter, r *http.Request) {
	overview, err := h.analyticsService.GetSystemOverview()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(overview)
}
