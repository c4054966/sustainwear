package analytics

import (
	"errors"
	"log"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GETS DONATION TRENDS
func (s *Service) GetDonationTrends(orgID uint, period string, startDate, endDate time.Time) ([]DonationTrend, error) {
	if orgID == 0 {
		return nil, errors.New("organisation ID is required")
	}

	validPeriods := map[string]bool{"daily": true, "weekly": true, "monthly": true}
	if !validPeriods[period] {
		return nil, errors.New("invalid period, must be daily, weekly, or monthly")
	}

	if endDate.Before(startDate) {
		return nil, errors.New("end date must be after start date")
	}

	trends, err := s.repo.GetDonationTrends(orgID, period, startDate, endDate)
	if err != nil {
		log.Printf("ANALYTICS: Failed to get donation trends: %v", err)
		return nil, err
	}

	log.Printf("ANALYTICS: Retrieved %d donation trends for org %d", len(trends), orgID)
	return trends, nil
}

// GETS CATEGORY BREAKDOWN
func (s *Service) GetCategoryBreakdown(orgID uint) ([]CategoryBreakdown, error) {
	if orgID == 0 {
		return nil, errors.New("organisation ID is required")
	}

	breakdown, err := s.repo.GetCategoryBreakdown(orgID)
	if err != nil {
		log.Printf("ANALYTICS: Failed to get category breakdown: %v", err)
		return nil, err
	}

	log.Printf("ANALYTICS: Retrieved category breakdown for org %d", orgID)
	return breakdown, nil
}

// GETS SUSTAINABILITY METRICS
func (s *Service) GetSustainabilityMetrics(orgID uint, period string) (*SustainabilityMetrics, error) {
	if orgID == 0 {
		return nil, errors.New("organisation ID is required")
	}

	validPeriods := map[string]bool{"all_time": true, "last_month": true, "last_year": true}
	if !validPeriods[period] {
		return nil, errors.New("invalid period, must be all_time, last_month, or last_year")
	}

	metrics, err := s.repo.GetSustainabilityMetrics(orgID, period)
	if err != nil {
		log.Printf("ANALYTICS: Failed to get sustainability metrics: %v", err)
		return nil, err
	}

	log.Printf("ANALYTICS: Retrieved sustainability metrics for org %d (%s)", orgID, period)
	return metrics, nil
}

// GETS DONOR IMPACT
func (s *Service) GetDonorImpact(donorID uint) (*DonorImpact, error) {
	if donorID == 0 {
		return nil, errors.New("donor ID is required")
	}

	impact, err := s.repo.GetDonorImpact(donorID)
	if err != nil {
		log.Printf("ANALYTICS: Failed to get donor impact: %v", err)
		return nil, err
	}

	log.Printf("ANALYTICS: Retrieved donor impact for donor %d", donorID)
	return impact, nil
}

// GETS ORG PERFORMANCE
func (s *Service) GetOrgPerformance(orgID uint) (*OrgPerformance, error) {
	if orgID == 0 {
		return nil, errors.New("organisation ID is required")
	}

	performance, err := s.repo.GetOrgPerformance(orgID)
	if err != nil {
		log.Printf("ANALYTICS: Failed to get org performance: %v", err)
		return nil, err
	}

	log.Printf("ANALYTICS: Retrieved performance metrics for org %d", orgID)
	return performance, nil
}

// GETS SYSTEM OVERVIEW (ADMIN ONLY)
func (s *Service) GetSystemOverview() (*SystemOverview, error) {
	overview, err := s.repo.GetSystemOverview()
	if err != nil {
		log.Printf("ANALYTICS: Failed to get system overview: %v", err)
		return nil, err
	}

	log.Printf("ANALYTICS: Retrieved system-wide overview")
	return overview, nil
}
