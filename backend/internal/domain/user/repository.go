package user

import (
	"database/sql"
	"errors"
	"log"
)

type Repository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	List() ([]User, error)
	ListPaginated(limit, offset int) ([]User, error)
	Update(user *User) error
	Delete(id uint) error
}

type SQLRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(user *User) error {
	query := `INSERT INTO users (email, password_hash, full_name, role, org_id, is_active, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`

	result, err := r.db.Exec(query, user.Email, user.PasswordHash, user.FullName, user.Role, user.OrganisationID, user.IsActive)
	if err != nil {
		log.Printf("USER: Failed to create user: %v", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("USER: Failed to get user ID: %v", err)
		return err
	}

	user.ID = uint(id)
	return nil
}

func (r *SQLRepository) GetByID(id uint) (*User, error) {
	var user User
	query := `SELECT id, email, password_hash, full_name, role, org_id, is_active, created_at, updated_at 
	          FROM users WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Role, &user.OrganisationID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		log.Printf("USER: User not found with ID: %d", id)
		return nil, errors.New("user not found")
	}
	if err != nil {
		log.Printf("USER: Failed to get user by ID: %v", err)
		return nil, err
	}

	return &user, nil
}

func (r *SQLRepository) GetByEmail(email string) (*User, error) {
	var user User
	query := `SELECT id, email, password_hash, full_name, role, org_id, is_active, created_at, updated_at 
	          FROM users WHERE email = ? AND is_active = 1`

	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Role, &user.OrganisationID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		log.Printf("USER: User not found with email: %s", email)
		return nil, errors.New("user not found")
	}
	if err != nil {
		log.Printf("USER: Failed to get user by email: %v", err)
		return nil, err
	}

	return &user, nil
}

func (r *SQLRepository) List() ([]User, error) {
	query := `SELECT id, email, full_name, role, org_id, is_active, created_at, updated_at 
	          FROM users ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("USER: Failed to list users: %v", err)
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Email, &u.FullName, &u.Role, &u.OrganisationID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			continue
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *SQLRepository) ListPaginated(limit, offset int) ([]User, error) {
	query := `SELECT id, email, full_name, role, org_id, is_active, created_at, updated_at
	FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		log.Printf("USER: Failed to list users: %v", err)
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Email, &u.FullName, &u.Role, &u.OrganisationID,
			&u.IsActive, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			continue
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *SQLRepository) Update(user *User) error {
	query := `UPDATE users SET full_name = ?, role = ?, org_id = ?, is_active = ?, updated_at = datetime('now') 
	          WHERE id = ?`

	result, err := r.db.Exec(query, user.FullName, user.Role, user.OrganisationID, user.IsActive, user.ID)
	if err != nil {
		log.Printf("USER: Failed to update user: %v", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("USER: Failed to get rows affected: %v", err)
		return err
	}

	if rows == 0 {
		log.Printf("USER: User not found with ID: %d", user.ID)
		return errors.New("user not found")
	}

	return nil
}

func (r *SQLRepository) Delete(id uint) error {
	query := `UPDATE users SET is_active = 0, updated_at = datetime('now') WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("USER: Failed to delete user: %v", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("USER: Failed to get rows affected: %v", err)
		return err
	}

	if rows == 0 {
		log.Printf("USER: User not found with ID: %d", id)
		return errors.New("user not found")
	}

	return nil
}
