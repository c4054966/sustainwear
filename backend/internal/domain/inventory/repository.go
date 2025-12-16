package inventory

import (
	"database/sql"
	"errors"
)

type Repository interface {
	Create(item *InventoryItem) error
	GetByID(id uint) (*InventoryItem, error)
	List(orgID uint, filters map[string]interface{}) ([]InventoryItem, error)
	Update(item *InventoryItem) error
	UpdateQuantities(id uint, available, allocated, distributed int) error
	Delete(id uint) error
	GetStats(orgID uint) (*InventoryStats, error)
}

type SQLRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &SQLRepository{db: db}
}

// CREATES INVENTORY ITEM
func (r *SQLRepository) Create(item *InventoryItem) error {
	query := `INSERT INTO inventory (donation_id, item_name, category, condition, quantity, available_qty, allocated_qty, distributed_qty, location, status, org_id, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`

	result, err := r.db.Exec(query, item.DonationID, item.ItemName, item.Category, item.Condition,
		item.Quantity, item.AvailableQty, item.AllocatedQty, item.DistributedQty,
		item.Location, item.Status, item.OrgID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	item.ID = uint(id)
	return nil
}

// GETS INVENTORY ITEM BY ID
func (r *SQLRepository) GetByID(id uint) (*InventoryItem, error) {
	var item InventoryItem
	query := `SELECT id, donation_id, item_name, category, condition, quantity, available_qty, allocated_qty, distributed_qty, location, status, org_id, created_at, updated_at
	          FROM inventory WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(&item.ID, &item.DonationID, &item.ItemName, &item.Category,
		&item.Condition, &item.Quantity, &item.AvailableQty, &item.AllocatedQty, &item.DistributedQty,
		&item.Location, &item.Status, &item.OrgID, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("inventory item not found")
	}

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// LISTS INVENTORY ITEMS
func (r *SQLRepository) List(orgID uint, filters map[string]interface{}) ([]InventoryItem, error) {
	query := `SELECT id, donation_id, item_name, category, condition, quantity, available_qty, allocated_qty, distributed_qty, location, status, org_id, created_at, updated_at
	          FROM inventory WHERE org_id = ?`

	args := []interface{}{orgID}

	if category, ok := filters["category"].(string); ok && category != "" {
		query += ` AND category = ?`
		args = append(args, category)
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}

	if condition, ok := filters["condition"].(string); ok && condition != "" {
		query += ` AND condition = ?`
		args = append(args, condition)
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []InventoryItem{}
	for rows.Next() {
		var item InventoryItem
		err := rows.Scan(&item.ID, &item.DonationID, &item.ItemName, &item.Category, &item.Condition,
			&item.Quantity, &item.AvailableQty, &item.AllocatedQty, &item.DistributedQty,
			&item.Location, &item.Status, &item.OrgID, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

// UPDATES INVENTORY ITEM
func (r *SQLRepository) Update(item *InventoryItem) error {
	query := `UPDATE inventory SET location = ?, status = ?, available_qty = ?, allocated_qty = ?, distributed_qty = ?, updated_at = datetime('now')
	          WHERE id = ?`

	result, err := r.db.Exec(query, item.Location, item.Status, item.AvailableQty, item.AllocatedQty, item.DistributedQty, item.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("inventory item not found")
	}

	return nil
}

// UPDATES INVENTORY QUANTITIES
func (r *SQLRepository) UpdateQuantities(id uint, available, allocated, distributed int) error {
	query := `UPDATE inventory SET available_qty = ?, allocated_qty = ?, distributed_qty = ?, updated_at = datetime('now')
	          WHERE id = ?`

	result, err := r.db.Exec(query, available, allocated, distributed, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("inventory item not found")
	}

	return nil
}

// DELETES INVENTORY ITEM
func (r *SQLRepository) Delete(id uint) error {
	query := `DELETE FROM inventory WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("inventory item not found")
	}

	return nil
}

// GETS INVENTORY STATISTICS
func (r *SQLRepository) GetStats(orgID uint) (*InventoryStats, error) {
	stats := &InventoryStats{
		ByCategory:  make(map[string]int),
		ByCondition: make(map[string]int),
	}

	// TOTAL ITEMS
	r.db.QueryRow(`SELECT COALESCE(SUM(quantity), 0) FROM inventory WHERE org_id = ?`, orgID).Scan(&stats.TotalItems)

	// AVAILABLE ITEMS
	r.db.QueryRow(`SELECT COALESCE(SUM(available_qty), 0) FROM inventory WHERE org_id = ?`, orgID).Scan(&stats.AvailableItems)

	// ALLOCATED ITEMS
	r.db.QueryRow(`SELECT COALESCE(SUM(allocated_qty), 0) FROM inventory WHERE org_id = ?`, orgID).Scan(&stats.AllocatedItems)

	// DISTRIBUTED ITEMS
	r.db.QueryRow(`SELECT COALESCE(SUM(distributed_qty), 0) FROM inventory WHERE org_id = ?`, orgID).Scan(&stats.DistributedItems)

	// VIA CATEGORY
	rows, _ := r.db.Query(`SELECT category, SUM(quantity) FROM inventory WHERE org_id = ? GROUP BY category`, orgID)
	defer rows.Close()
	for rows.Next() {
		var category string
		var count int
		rows.Scan(&category, &count)
		stats.ByCategory[category] = count
	}

	// VIA CONDITION
	rows2, _ := r.db.Query(`SELECT condition, SUM(quantity) FROM inventory WHERE org_id = ? GROUP BY condition`, orgID)
	defer rows2.Close()
	for rows2.Next() {
		var condition string
		var count int
		rows2.Scan(&condition, &count)
		stats.ByCondition[condition] = count
	}

	// LOW STOCK ALERTS (ITEMS WITH available_qty < 5)
	rows3, _ := r.db.Query(`SELECT item_name FROM inventory WHERE org_id = ? AND available_qty < 5 AND available_qty > 0`, orgID)
	defer rows3.Close()
	lowStock := []string{}
	for rows3.Next() {
		var itemName string
		rows3.Scan(&itemName)
		lowStock = append(lowStock, itemName)
	}
	stats.LowStockAlerts = lowStock

	return stats, nil
}
