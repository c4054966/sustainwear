package inventory

import "time"

type InventoryItem struct {
	ID             uint      `json:"id" db:"id"`
	DonationID     uint      `json:"donation_id" db:"donation_id"`
	ItemName       string    `json:"item_name" db:"item_name"`
	Category       string    `json:"category" db:"category"`
	Condition      string    `json:"condition" db:"condition"`
	Quantity       int       `json:"quantity" db:"quantity"`
	AvailableQty   int       `json:"available_qty" db:"available_qty"`
	AllocatedQty   int       `json:"allocated_qty" db:"allocated_qty"`
	DistributedQty int       `json:"distributed_qty" db:"distributed_qty"`
	Location       string    `json:"location" db:"location"`
	Status         string    `json:"status" db:"status"` // available, allocated, distributed
	OrgID          uint      `json:"org_id" db:"org_id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CreateInventoryRequest struct {
	ItemName    string `json:"item_name"`
	Category    string `json:"category"`
	Condition   string `json:"condition"`
	Quantity    int    `json:"quantity"`
	Location    string `json:"location"`
	Description string `json:"description,omitempty"`
}

type UpdateInventoryRequest struct {
	Quantity *int    `json:"quantity"`
	Location *string `json:"location"`
	Status   *string `json:"status"`
}

type DistributeRequest struct {
	RecipientInfo       string `json:"recipient_info"`
	QuantityDistributed int    `json:"quantity_distributed"`
	DistributionDate    string `json:"distribution_date"`
}

type InventoryStats struct {
	TotalItems       int            `json:"total_items"`
	AvailableItems   int            `json:"available_items"`
	AllocatedItems   int            `json:"allocated_items"`
	DistributedItems int            `json:"distributed_items"`
	ByCategory       map[string]int `json:"by_category"`
	ByCondition      map[string]int `json:"by_condition"`
	LowStockAlerts   []string       `json:"low_stock_alerts"`
}
