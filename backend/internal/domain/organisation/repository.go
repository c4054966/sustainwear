package organisation

import (
	"database/sql"
	"errors"
	"log"
)

type Repository interface {
	Create(org *Organisation) error
	GetByID(id uint) (*Organisation, error)
	GetByEmail(email string) (*Organisation, error)
	List(filters map[string]interface{}) ([]Organisation, error)
	Update(org *Organisation) error
	Delete(id uint) error
	GetStats(orgID uint) (*OrgStats, error)
}

type SQLRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &SQLRepository{db: db}
}

// CREATES ORGANISATION
func (r *SQLRepository) Create(org *Organisation) error {
	query := `INSERT INTO organisations (name, description, type, email, phone, address, city, county, postcode, country, website, status, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`

	result, err := r.db.Exec(query, org.Name, org.Description, org.Type, org.Email, org.Phone,
		org.Address, org.City, org.County, org.PostCode, org.Country, org.Website, org.Status)
	if err != nil {
		log.Printf("ORGANISATION: Failed to create organisation: %v", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("ORGANISATION: Failed to get organisation ID: %v", err)
		return err
	}

	org.ID = uint(id)
	return nil
}

// GETS ORGANISATION BY ID
func (r *SQLRepository) GetByID(id uint) (*Organisation, error) {
	var org Organisation
	query := `SELECT id, name, description, type, email, phone, address, city, county, postcode, country, website, status, created_at, updated_at
	          FROM organisations WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(&org.ID, &org.Name, &org.Description, &org.Type, &org.Email,
		&org.Phone, &org.Address, &org.City, &org.County, &org.PostCode, &org.Country, &org.Website,
		&org.Status, &org.CreatedAt, &org.UpdatedAt)

	if err == sql.ErrNoRows {
		log.Printf("ORGANISATION: Organisation not found with ID: %d", id)
		return nil, errors.New("organisation not found")
	}

	if err != nil {
		log.Printf("ORGANISATION: Failed to get organisation by ID: %v", err)
		return nil, err
	}

	return &org, nil
}

// GETS ORGANISATION BY EMAIL
func (r *SQLRepository) GetByEmail(email string) (*Organisation, error) {
	var org Organisation
	query := `SELECT id, name, description, type, email, phone, address, city, county, postcode, country, website, status, created_at, updated_at
	          FROM organisations WHERE email = ?`

	err := r.db.QueryRow(query, email).Scan(&org.ID, &org.Name, &org.Description, &org.Type, &org.Email,
		&org.Phone, &org.Address, &org.City, &org.County, &org.PostCode, &org.Country, &org.Website,
		&org.Status, &org.CreatedAt, &org.UpdatedAt)

	if err == sql.ErrNoRows {
		log.Printf("ORGANISATION: Organisation not found with email: %s", email)
		return nil, errors.New("organisation not found")
	}

	if err != nil {
		log.Printf("ORGANISATION: Failed to get organisation by email: %v", err)
		return nil, err
	}

	return &org, nil
}

// LISTS ORGANISATIONS
func (r *SQLRepository) List(filters map[string]interface{}) ([]Organisation, error) {
	query := `SELECT id, name, description, type, email, phone, address, city, county, postcode, country, website, status, created_at, updated_at
	          FROM organisations WHERE 1=1`

	args := []interface{}{}

	if orgType, ok := filters["type"].(string); ok && orgType != "" {
		query += ` AND type = ?`
		args = append(args, orgType)
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}

	if city, ok := filters["city"].(string); ok && city != "" {
		query += ` AND city = ?`
		args = append(args, city)
	}

	if county, ok := filters["county"].(string); ok && county != "" {
		query += ` AND county = ?`
		args = append(args, county)
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Printf("ORGANISATION: Failed to list organisations: %v", err)
		return nil, err
	}
	defer rows.Close()

	organisations := []Organisation{}
	for rows.Next() {
		var org Organisation
		err := rows.Scan(&org.ID, &org.Name, &org.Description, &org.Type, &org.Email, &org.Phone,
			&org.Address, &org.City, &org.County, &org.PostCode, &org.Country, &org.Website,
			&org.Status, &org.CreatedAt, &org.UpdatedAt)
		if err != nil {
			continue
		}
		organisations = append(organisations, org)
	}

	return organisations, nil
}

// UPDATES ORGANISATION
func (r *SQLRepository) Update(org *Organisation) error {
	query := `UPDATE organisations SET name = ?, description = ?, phone = ?, address = ?, city = ?, county = ?, postcode = ?, website = ?, status = ?, updated_at = datetime('now')
	          WHERE id = ?`

	result, err := r.db.Exec(query, org.Name, org.Description, org.Phone, org.Address, org.City,
		org.County, org.PostCode, org.Website, org.Status, org.ID)
	if err != nil {
		log.Printf("ORGANISATION: Failed to update organisation: %v", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("ORGANISATION: Failed to get rows affected: %v", err)
		return err
	}

	if rows == 0 {
		log.Printf("ORGANISATION: Organisation not found with ID: %d", org.ID)
		return errors.New("organisation not found")
	}

	return nil
}

// DELETES ORGANISATION
func (r *SQLRepository) Delete(id uint) error {
	query := `DELETE FROM organisations WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("ORGANISATION: Failed to delete organisation: %v", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("ORGANISATION: Failed to get rows affected: %v", err)
		return err
	}

	if rows == 0 {
		log.Printf("ORGANISATION: Organisation not found with ID: %d", id)
		return errors.New("organisation not found")
	}

	return nil
}

// GETS ORGANISATION STATISTICS (UK CLOTHES DONATION FOCUSED)
func (r *SQLRepository) GetStats(orgID uint) (*OrgStats, error) {
	stats := &OrgStats{}

	// TOTAL DONATIONS RECEIVED
	r.db.QueryRow(`SELECT COUNT(*) FROM donations WHERE org_id = ?`, orgID).Scan(&stats.TotalDonations)

	// TOTAL INVENTORY QUANTITY
	r.db.QueryRow(`SELECT COALESCE(SUM(quantity), 0) FROM inventory WHERE org_id = ?`, orgID).Scan(&stats.TotalInventoryQty)

	// AVAILABLE STOCK (ITEMS READY TO ALLOCATE)
	r.db.QueryRow(`SELECT COALESCE(SUM(available_qty), 0) FROM inventory WHERE org_id = ? AND status = 'available'`, orgID).Scan(&stats.AvailableStock)

	// ACTIVE CHARITY STAFF (NOT DONORS)
	r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE org_id = ? AND role IN ('charity_staff', 'admin')`, orgID).Scan(&stats.ActiveStaff)

	// SUSTAINABILITY METRICS - CO2 SAVED (ESTIMATE: 6KG CO2 PER KG OF CLOTHES DONATED)
	var totalWeight float64
	r.db.QueryRow(`SELECT COALESCE(SUM(quantity), 0) FROM inventory WHERE org_id = ?`, orgID).Scan(&totalWeight)
	stats.CO2SavedKg = totalWeight * 6.0 // Rough estimate

	// LANDFILL REDUCTION (ASSUME AVERAGE CLOTHING WEIGHT 0.5KG PER ITEM)
	stats.LandfillReductionKg = totalWeight * 0.5

	return stats, nil
}
