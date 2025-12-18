package donation

import (
	"database/sql"
	"errors"

	jsoniter "github.com/json-iterator/go"
)

type Repository interface {
	Create(donation *Donation) error
	GetByID(id uint) (*Donation, error)
	List(filters map[string]interface{}) ([]Donation, error)
	UpdateStatus(id uint, status, notes string) error
	Delete(id uint) error
}

type SQLRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &SQLRepository{db: db}
}

// CREATES DONATION
func (r *SQLRepository) Create(donation *Donation) error {
	query := `INSERT INTO donations (donor_id, item_name, description, category, size, gender, condition, quantity, images, status, org_id, notes, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`

	result, err := r.db.Exec(query, donation.DonorID, donation.ItemName, donation.Description,
		donation.Category, donation.Size, donation.Gender, donation.Condition,
		donation.Quantity, donation.Images, donation.Status, donation.OrgID, donation.Notes)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	donation.ID = uint(id)
	return nil
}

// GETS DONATION BY ID
func (r *SQLRepository) GetByID(id uint) (*Donation, error) {
	var donation Donation
	query := `SELECT id, donor_id, item_name, description, category, size, gender, condition, quantity, images, status, org_id, notes, created_at, updated_at
	          FROM donations WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(&donation.ID, &donation.DonorID, &donation.ItemName,
		&donation.Description, &donation.Category, &donation.Size, &donation.Gender,
		&donation.Condition, &donation.Quantity, &donation.Images, &donation.Status,
		&donation.OrgID, &donation.Notes, &donation.CreatedAt, &donation.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("donation not found")
	}

	if err != nil {
		return nil, err
	}

	return &donation, nil
}

// LISTS DONATIONS
func (r *SQLRepository) List(filters map[string]interface{}) ([]Donation, error) {
	query := `SELECT id, donor_id, item_name, description, category, size, gender, condition, quantity, images, status, org_id, created_at, updated_at
	          FROM donations WHERE 1=1`

	args := []interface{}{}

	if donorID, ok := filters["donor_id"].(uint); ok && donorID > 0 {
		query += ` AND donor_id = ?`
		args = append(args, donorID)
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}

	if category, ok := filters["category"].(string); ok && category != "" {
		query += ` AND category = ?`
		args = append(args, category)
	}

	if orgID, ok := filters["org_id"].(uint); ok && orgID > 0 {
		query += ` AND org_id = ?`
		args = append(args, orgID)
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	donations := []Donation{}
	for rows.Next() {
		var d Donation
		err := rows.Scan(&d.ID, &d.DonorID, &d.ItemName, &d.Description, &d.Category,
			&d.Size, &d.Gender, &d.Condition, &d.Quantity, &d.Images, &d.Status,
			&d.OrgID, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			continue
		}
		donations = append(donations, d)
	}

	return donations, nil
}

// UPDATES DONATION STATUS
func (r *SQLRepository) UpdateStatus(id uint, status, notes string) error {
	query := `UPDATE donations SET status = ?, notes = ?, updated_at = datetime('now') WHERE id = ?`

	result, err := r.db.Exec(query, status, notes, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("donation not found")
	}

	return nil
}

// DELETES DONATION BY ID
func (r *SQLRepository) Delete(id uint) error {
	query := `DELETE FROM donations WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("donation not found")
	}

	return nil
}

// MARSHALS IMAGE ARRAY TO JSON STRING
func MarshalImages(images []string) string {
	if len(images) == 0 {
		return "[]"
	}
	data, _ := jsoniter.Marshal(images)
	return string(data)
}

// UNMARSHALS JSON STRING TO IMAGE ARRAY
func UnmarshalImages(imagesJSON string) []string {
	var images []string
	jsoniter.Unmarshal([]byte(imagesJSON), &images)
	return images
}
