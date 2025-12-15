package analytics

import (
	"database/sql"
	"log"
	"strings"
	"time"
)

type Repository interface {
	GetDonationTrends(orgID uint, period string, startDate, endDate time.Time) ([]DonationTrend, error)
	GetCategoryBreakdown(orgID uint) ([]CategoryBreakdown, error)
	GetSustainabilityMetrics(orgID uint, period string) (*SustainabilityMetrics, error)
	GetDonorImpact(donorID uint) (*DonorImpact, error)
	GetOrgPerformance(orgID uint) (*OrgPerformance, error)
	GetSystemOverview() (*SystemOverview, error)
}

type SQLRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &SQLRepository{db: db}
}

// GETS DONATION TRENDS BY PERIOD
func (r *SQLRepository) GetDonationTrends(orgID uint, period string, startDate, endDate time.Time) ([]DonationTrend, error) {
	var query string
	var dateFormat string

	switch strings.ToLower(period) {
	case "daily":
		dateFormat = "%d-%m-%Y"
	case "weekly":
		dateFormat = "%Y-W%W"
	case "monthly":
		dateFormat = "%m-%Y"
	default:
		dateFormat = "%d-%m-%Y"
	}

	query = `SELECT strftime(?, created_at) as period, COUNT(*) as total_donations, COALESCE(SUM(quantity), 0) as total_items
	         FROM donations 
	         WHERE org_id = ? AND created_at BETWEEN ? AND ?
	         GROUP BY period
	         ORDER BY period ASC`

	rows, err := r.db.Query(query, dateFormat, orgID, startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Printf("ANALYTICS: Failed to get donation trends: %v", err)
		return nil, err
	}
	defer rows.Close()

	trends := []DonationTrend{}
	for rows.Next() {
		var trend DonationTrend
		err := rows.Scan(&trend.Timestamp, &trend.TotalDonations, &trend.TotalItems)
		if err != nil {
			continue
		}
		trend.Period = period
		trends = append(trends, trend)
	}

	return trends, nil
}

// GETS CATEGORY BREAKDOWN FOR ORG
func (r *SQLRepository) GetCategoryBreakdown(orgID uint) ([]CategoryBreakdown, error) {
	query := `SELECT category, COUNT(*) as count
	          FROM inventory
	          WHERE org_id = ?
	          GROUP BY category
	          ORDER BY count DESC`

	rows, err := r.db.Query(query, orgID)
	if err != nil {
		log.Printf("ANALYTICS: Failed to get category breakdown: %v", err)
		return nil, err
	}
	defer rows.Close()

	breakdown := []CategoryBreakdown{}
	for rows.Next() {
		var cb CategoryBreakdown
		err := rows.Scan(&cb.Category, &cb.Count)
		if err != nil {
			continue
		}
		cb.OrgID = orgID
		breakdown = append(breakdown, cb)
	}

	return breakdown, nil
}

// GETS SUSTAINABILITY METRICS
func (r *SQLRepository) GetSustainabilityMetrics(orgID uint, period string) (*SustainabilityMetrics, error) {
	metrics := &SustainabilityMetrics{Period: period}

	var dateFilter string
	switch period {
	case "last_month":
		dateFilter = "AND created_at >= datetime('now', '-1 month')"
	case "last_year":
		dateFilter = "AND created_at >= datetime('now', '-1 year')"
	default:
		dateFilter = ""
	}

	query := `SELECT COUNT(*) FROM donations WHERE org_id = ? ` + dateFilter
	r.db.QueryRow(query, orgID).Scan(&metrics.TotalDonations)

	query = `SELECT COALESCE(SUM(quantity), 0) FROM donations WHERE org_id = ? ` + dateFilter
	r.db.QueryRow(query, orgID).Scan(&metrics.TotalItemsDonated)

	// CO2 CALCULATION (6KG PER ITEM)
	metrics.CO2SavedKg = float64(metrics.TotalItemsDonated) * 6.0

	// LANDFILL REDUCTION (0.5KG PER ITEM)
	metrics.LandfillReductionKg = float64(metrics.TotalItemsDonated) * 0.5

	// BENEFICIARIES HELPED (5 ITEMS PER PERSON)
	var distributed int
	query = `SELECT COALESCE(SUM(distributed_qty), 0) FROM inventory WHERE org_id = ?`
	r.db.QueryRow(query, orgID).Scan(&distributed)
	metrics.BeneficiariesHelped = distributed / 5

	return metrics, nil
}

// GETS DONOR IMPACT STATS
func (r *SQLRepository) GetDonorImpact(donorID uint) (*DonorImpact, error) {
	impact := &DonorImpact{DonorID: donorID}

	query := `SELECT COUNT(*), COALESCE(SUM(quantity), 0), MIN(created_at), MAX(created_at)
	          FROM donations 
	          WHERE donor_id = ?`

	var firstDonation, lastDonation string
	err := r.db.QueryRow(query, donorID).Scan(&impact.TotalDonations, &impact.TotalItemsDonated, &firstDonation, &lastDonation)
	if err != nil {
		log.Printf("ANALYTICS: Failed to get donor impact: %v", err)
		return nil, err
	}

	// UK DATE FORMAT PARSING
	impact.FirstDonation, _ = time.Parse("2006-01-02 15:04:05", firstDonation)
	impact.LastDonation, _ = time.Parse("2006-01-02 15:04:05", lastDonation)

	// CO2 CALCULATION
	impact.CO2SavedKg = float64(impact.TotalItemsDonated) * 6.0

	// LANDFILL REDUCTION
	impact.LandfillReductionKg = float64(impact.TotalItemsDonated) * 0.5

	return impact, nil
}

// GETS ORG PERFORMANCE METRICS
func (r *SQLRepository) GetOrgPerformance(orgID uint) (*OrgPerformance, error) {
	perf := &OrgPerformance{OrgID: orgID}

	r.db.QueryRow(`SELECT name FROM organisations WHERE id = ?`, orgID).Scan(&perf.OrgName)

	r.db.QueryRow(`SELECT COUNT(*) FROM donations WHERE org_id = ?`, orgID).Scan(&perf.DonationsReceived)

	r.db.QueryRow(`SELECT COALESCE(SUM(quantity), 0) FROM inventory WHERE org_id = ?`, orgID).Scan(&perf.ItemsProcessed)

	r.db.QueryRow(`SELECT COALESCE(SUM(available_qty), 0) FROM inventory WHERE org_id = ?`, orgID).Scan(&perf.ItemsAvailable)

	r.db.QueryRow(`SELECT COALESCE(SUM(allocated_qty), 0) FROM inventory WHERE org_id = ?`, orgID).Scan(&perf.ItemsAllocated)

	r.db.QueryRow(`SELECT COALESCE(SUM(distributed_qty), 0) FROM inventory WHERE org_id = ?`, orgID).Scan(&perf.ItemsDistributed)

	// UTILIZATION RATE CALCULATION
	if perf.ItemsProcessed > 0 {
		perf.UtilizationRate = float64(perf.ItemsDistributed) / float64(perf.ItemsProcessed) * 100
	}

	return perf, nil
}

// GETS SYSTEM-WIDE OVERVIEW
func (r *SQLRepository) GetSystemOverview() (*SystemOverview, error) {
	overview := &SystemOverview{}

	r.db.QueryRow(`SELECT COUNT(*) FROM organisations WHERE status = 'active'`).Scan(&overview.TotalOrganisations)

	r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE role = 'donor'`).Scan(&overview.TotalDonors)

	r.db.QueryRow(`SELECT COUNT(*) FROM donations`).Scan(&overview.TotalDonations)

	r.db.QueryRow(`SELECT COALESCE(SUM(quantity), 0) FROM inventory`).Scan(&overview.TotalItemsProcessed)

	// CO2 CALCULATION
	overview.CO2SavedKg = float64(overview.TotalItemsProcessed) * 6.0

	// LANDFILL REDUCTION
	overview.LandfillReductionKg = float64(overview.TotalItemsProcessed) * 0.5

	return overview, nil
}
