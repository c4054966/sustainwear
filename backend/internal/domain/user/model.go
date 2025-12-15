package user

import "time"

const (
	RoleDonor   = "donor"
	RoleCharity = "charity_staff"
	RoleAdmin   = "admin"
)

type User struct {
	ID             uint      `json:"id" db:"id"`
	Email          string    `json:"email" db:"email"`
	PasswordHash   string    `json:"-" db:"password_hash"`
	FullName       string    `json:"full_name" db:"full_name"`
	Role           string    `json:"role" db:"role"`
	OrganisationID *uint     `json:"org_id,omitempty" db:"org_id"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	FullName       string `json:"full_name"`
	Role           string `json:"role"`
	OrganisationID *uint  `json:"org_id,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token    string `json:"token"`
	UserID   uint   `json:"user_id"`
	Role     string `json:"role"`
	FullName string `json:"full_name"`
}
