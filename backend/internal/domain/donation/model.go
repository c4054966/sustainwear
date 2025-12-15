package donation

import "time"

const (
	StatusPending     = "pending"
	StatusApproved    = "approved"
	StatusRejected    = "rejected"
	StatusDistributed = "distributed"
)

type Donation struct {
	ID          uint      `json:"id" db:"id"`
	DonorID     uint      `json:"donor_id" db:"donor_id"`
	ItemName    string    `json:"item_name" db:"item_name"`
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"`
	Size        string    `json:"size" db:"size"`
	Gender      string    `json:"gender" db:"gender"`
	Condition   string    `json:"condition" db:"condition"`
	Quantity    int       `json:"quantity" db:"quantity"`
	Images      string    `json:"images" db:"images"` // JSON STRING ARRAY SINCE SQLITE DOESN'T HAVE NATIVE ARRAY SUPPORT
	Status      string    `json:"status" db:"status"`
	OrgID       *uint     `json:"org_id,omitempty" db:"org_id"`
	Notes       string    `json:"notes,omitempty" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateRequest struct {
	ItemName    string   `json:"item_name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Size        string   `json:"size"`
	Gender      string   `json:"gender"`
	Condition   string   `json:"condition"`
	Quantity    int      `json:"quantity"`
	Images      []string `json:"images"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
	Notes  string `json:"notes"`
}
