package analytics

import "time"

type DonationTrend struct {
	Period         string `json:"period"` // daily, weekly, monthly
	TotalDonations int    `json:"total_donations"`
	TotalItems     int    `json:"total_items"`
	Timestamp      string `json:"timestamp"`
}

type CategoryBreakdown struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
	OrgID    uint   `json:"org_id"`
}

type SustainabilityMetrics struct {
	TotalDonations      int     `json:"total_donations"`
	TotalItemsDonated   int     `json:"total_items_donated"`
	CO2SavedKg          float64 `json:"co2_saved_kg"`
	LandfillReductionKg float64 `json:"landfill_reduction_kg"`
	BeneficiariesHelped int     `json:"beneficiaries_helped"`
	Period              string  `json:"period"` // all_time, last_month, last_year
}

type DonorImpact struct {
	DonorID             uint      `json:"donor_id"`
	TotalDonations      int       `json:"total_donations"`
	TotalItemsDonated   int       `json:"total_items_donated"`
	CO2SavedKg          float64   `json:"co2_saved_kg"`
	LandfillReductionKg float64   `json:"landfill_reduction_kg"`
	FirstDonation       time.Time `json:"first_donation"`
	LastDonation        time.Time `json:"last_donation"`
}

type OrgPerformance struct {
	OrgID             uint    `json:"org_id"`
	OrgName           string  `json:"organization_name"`
	DonationsReceived int     `json:"donations_received"`
	ItemsProcessed    int     `json:"items_processed"`
	ItemsAvailable    int     `json:"items_available"`
	ItemsAllocated    int     `json:"items_allocated"`
	ItemsDistributed  int     `json:"items_distributed"`
	UtilizationRate   float64 `json:"utilization_rate"` // DISTRIBUTED / TOTAL RECEIVED
}

type SystemOverview struct {
	TotalOrganisations  int     `json:"total_organisations"`
	TotalDonors         int     `json:"total_donors"`
	TotalDonations      int     `json:"total_donations"`
	TotalItemsProcessed int     `json:"total_items_processed"`
	CO2SavedKg          float64 `json:"co2_saved_kg"`
	LandfillReductionKg float64 `json:"landfill_reduction_kg"`
}
