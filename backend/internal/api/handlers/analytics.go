package handlers

import (
	"log"
	"net/http"
	"time"

	"sustainwear/internal/api/middleware"
	"sustainwear/internal/config"
	"sustainwear/internal/domain/analytics"
	"sustainwear/internal/domain/user"

	jsoniter "github.com/json-iterator/go"
)

type AnalyticsHandler struct {
	analyticsService *analytics.Service
	userService      *user.Service
	config           *config.Config
}

func NewAnalyticsHandler(analyticsService *analytics.Service, userService *user.Service, cfg *config.Config) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		userService:      userService,
		config:           cfg,
	}
}

// GET DONATION TRENDS
func (h *AnalyticsHandler) GetDonationTrends(w http.ResponseWriter, r *http.Request) {
	orgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID)
	if err != nil {
		log.Printf("ANALYTICS: [GET api/analytics/trends] - Unauthorized access attempt - %v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
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
		log.Printf("ANALYTICS: [GET api/analytics/trends] - Invalid start_date format - %v", err)
		http.Error(w, "Invalid start_date format (use DD-MM-YYYY)", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("02-01-2006", endDateStr)
	if err != nil {
		log.Printf("ANALYTICS: [GET api/analytics/trends] - Invalid end_date format - %v", err)
		http.Error(w, "Invalid end_date format (use DD-MM-YYYY)", http.StatusBadRequest)
		return
	}

	trends, err := h.analyticsService.GetDonationTrends(uint(orgID), period, startDate, endDate)
	if err != nil {
		log.Printf("ANALYTICS: [GET api/analytics/trends] - Failed to get donation trends - %v", err)
		http.Error(w, "Unable to get donation trends", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(trends)
}

// GET CATEGORY BREAKDOWN
func (h *AnalyticsHandler) GetCategoryBreakdown(w http.ResponseWriter, r *http.Request) {
	orgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID)
	if err != nil {
		log.Printf("ANALYTICS: [GET api/analytics/categories] - Unauthorized access attempt - %v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	breakdown, err := h.analyticsService.GetCategoryBreakdown(uint(orgID))
	if err != nil {
		log.Printf("ANALYTICS: [GET api/analytics/categories] - Failed to get category breakdown - %v", err)
		http.Error(w, "Unable to get category breakdown", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(breakdown)
}

// GET SUSTAINABILITY METRICS
func (h *AnalyticsHandler) GetSustainabilityMetrics(w http.ResponseWriter, r *http.Request) {
	orgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID)
	if err != nil {
		log.Printf("ANALYTICS: [GET api/analytics/sustainability] - Unauthorized access attempt - %v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "all_time"
	}

	metrics, err := h.analyticsService.GetSustainabilityMetrics(uint(orgID), period)
	if err != nil {
		log.Printf("ANALYTICS: [GET api/analytics/sustainability] - Failed to get sustainability metrics - %v", err)
		http.Error(w, "Unable to get sustainability metrics", http.StatusInternalServerError)
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
		log.Printf("ANALYTICS: [GET api/analytics/donor-impact] - Failed to get donor impact - %v", err)
		http.Error(w, "Unable to get donor impact", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(impact)
}

// GET ORG PERFORMANCE
func (h *AnalyticsHandler) GetOrgPerformance(w http.ResponseWriter, r *http.Request) {
	orgID, err := middleware.GetOrgIDForRequest(r, h.userService.GetByID)
	if err != nil {
		log.Printf("ANALYTICS: [GET api/analytics/org-performance] - Unauthorized access attempt - %v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	performance, err := h.analyticsService.GetOrgPerformance(uint(orgID))
	if err != nil {
		log.Printf("ANALYTICS: [GET api/analytics/org-performance] - Failed to get organisation performance - %v", err)
		http.Error(w, "Unable to get organisation performance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(performance)
}

// GET SYSTEM OVERVIEW (ADMIN ONLY)
func (h *AnalyticsHandler) GetSystemOverview(w http.ResponseWriter, r *http.Request) {
	overview, err := h.analyticsService.GetSystemOverview()
	if err != nil {
		log.Printf("ANALYTICS: [GET api/analytics/system-overview] - Failed to get system overview - %v", err)
		http.Error(w, "Unable to get system overview", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(overview)
}
