package organisation

import "time"

type Organisation struct {
	ID          uint      `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Type        string    `json:"type" db:"type"` // charity, ngo, community, religious
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	Address     string    `json:"address" db:"address"`
	City        string    `json:"city" db:"city"`
	County      string    `json:"county" db:"county"`
	PostCode    string    `json:"postcode" db:"postcode"`
	Country     string    `json:"country" db:"country"`
	Website     string    `json:"website" db:"website"`
	Status      string    `json:"status" db:"status"` // active, inactive, pending
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateOrgRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	City        string `json:"city"`
	County      string `json:"county"`
	PostCode    string `json:"postcode"`
	Country     string `json:"country"`
	Website     string `json:"website"`
}

type UpdateOrgRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Phone       *string `json:"phone"`
	Address     *string `json:"address"`
	City        *string `json:"city"`
	County      *string `json:"county"`
	PostCode    *string `json:"postcode"`
	Website     *string `json:"website"`
	Status      *string `json:"status"`
}

type OrgStats struct {
	TotalDonations      int     `json:"total_donations"`
	TotalInventoryQty   int     `json:"total_inventory_quantity"`
	AvailableStock      int     `json:"available_stock"`
	ActiveStaff         int     `json:"active_staff"`
	CO2SavedKg          float64 `json:"co2_saved_kg"`          // SUSTAINABILITY
	LandfillReductionKg float64 `json:"landfill_reduction_kg"` // SUSTAINABILITY
}
