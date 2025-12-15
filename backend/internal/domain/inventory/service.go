package inventory

import (
	"errors"
	"log"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CREATES INVENTORY ITEM FROM DONATION
func (s *Service) CreateFromDonation(donationID uint, itemName, category, condition string, quantity int, orgID uint) (*InventoryItem, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	if itemName == "" || category == "" || condition == "" {
		return nil, errors.New("item name, category, and condition are required")
	}

	item := &InventoryItem{
		DonationID:     donationID,
		ItemName:       itemName,
		Category:       category,
		Condition:      condition,
		Quantity:       quantity,
		AvailableQty:   quantity, // ALL AVAILABLE INITIALLY
		AllocatedQty:   0,
		DistributedQty: 0,
		Location:       "warehouse", // DEFAULT LOCATION
		Status:         "available",
		OrgID:          orgID,
	}

	err := s.repo.Create(item)
	if err != nil {
		log.Printf("INVENTORY: Failed to create inventory item: %v", err)
		return nil, err
	}

	log.Printf("INVENTORY: Created inventory item ID %d for donation %d", item.ID, donationID)
	return item, nil
}

// GETS INVENTORY ITEM BY ID
func (s *Service) GetByID(id uint, orgID uint) (*InventoryItem, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// AUTH CHECK
	if item.OrgID != orgID {
		log.Printf("INVENTORY: Unauthorized access to inventory item %d by org %d", id, orgID)
		return nil, errors.New("unauthorized access")
	}

	return item, nil
}

// GETS INVENTORY ITEM BY DONATION ID
func (s *Service) GetByDonationID(donationID uint, orgID uint) (*InventoryItem, error) {
	item, err := s.repo.GetByDonationID(donationID)
	if err != nil {
		return nil, err
	}

	// AUTH CHECK
	if item.OrgID != orgID {
		log.Printf("INVENTORY: Unauthorized access to inventory by org %d", orgID)
		return nil, errors.New("unauthorized access")
	}

	return item, nil
}

// LISTS INVENTORY ITEMS
func (s *Service) List(orgID uint, filters map[string]interface{}) ([]InventoryItem, error) {
	items, err := s.repo.List(orgID, filters)
	if err != nil {
		log.Printf("INVENTORY: Failed to list inventory items: %v", err)
		return nil, err
	}

	log.Printf("INVENTORY: Listed %d inventory items for org %d", len(items), orgID)
	return items, nil
}

// UPDATES INVENTORY ITEM
func (s *Service) Update(id uint, orgID uint, req *UpdateInventoryRequest) (*InventoryItem, error) {
	item, err := s.GetByID(id, orgID)
	if err != nil {
		return nil, err
	}

	if req.Location != nil {
		item.Location = *req.Location
	}

	if req.Status != nil {
		validStatuses := map[string]bool{"available": true, "allocated": true, "distributed": true}
		if !validStatuses[*req.Status] {
			return nil, errors.New("invalid status")
		}
		item.Status = *req.Status
	}

	err = s.repo.Update(item)
	if err != nil {
		log.Printf("INVENTORY: Failed to update inventory item: %v", err)
		return nil, err
	}

	log.Printf("INVENTORY: Updated inventory item %d", id)
	return item, nil
}

// ALLOCATES INVENTORY FOR DISTRIBUTION
func (s *Service) Allocate(id uint, orgID uint, quantityToAllocate int) error {
	item, err := s.GetByID(id, orgID)
	if err != nil {
		return err
	}

	if quantityToAllocate <= 0 {
		return errors.New("allocation quantity must be greater than 0")
	}

	if item.AvailableQty < quantityToAllocate {
		return errors.New("insufficient available quantity")
	}

	// UPDATE QUANTITIES
	newAvailable := item.AvailableQty - quantityToAllocate
	newAllocated := item.AllocatedQty + quantityToAllocate

	err = s.repo.UpdateQuantities(id, newAvailable, newAllocated, item.DistributedQty)
	if err != nil {
		log.Printf("INVENTORY: Failed to allocate inventory: %v", err)
		return err
	}

	log.Printf("INVENTORY: Allocated %d units from inventory item %d", quantityToAllocate, id)
	return nil
}

// DISTRIBUTES ALLOCATED INVENTORY
func (s *Service) Distribute(id uint, orgID uint, quantityToDistribute int) error {
	item, err := s.GetByID(id, orgID)
	if err != nil {
		return err
	}

	if quantityToDistribute <= 0 {
		return errors.New("distribution quantity must be greater than 0")
	}

	if item.AllocatedQty < quantityToDistribute {
		return errors.New("insufficient allocated quantity")
	}

	// UPDATE QUANTITIES
	newAllocated := item.AllocatedQty - quantityToDistribute
	newDistributed := item.DistributedQty + quantityToDistribute

	err = s.repo.UpdateQuantities(id, item.AvailableQty, newAllocated, newDistributed)
	if err != nil {
		log.Printf("INVENTORY: Failed to distribute inventory: %v", err)
		return err
	}

	// UPDATE STATUS IF FULLY DISTRIBUTED
	if item.AvailableQty == 0 && newAllocated == 0 {
		item.Status = "distributed"
		s.repo.Update(item)
	}

	log.Printf("INVENTORY: Distributed %d units from inventory item %d", quantityToDistribute, id)
	return nil
}

// DEALLOCATES INVENTORY (RETURNS TO AVAILABLE)
func (s *Service) Deallocate(id uint, orgID uint, quantityToDeallocate int) error {
	item, err := s.GetByID(id, orgID)
	if err != nil {
		return err
	}

	if quantityToDeallocate <= 0 {
		return errors.New("deallocation quantity must be greater than 0")
	}

	if item.AllocatedQty < quantityToDeallocate {
		return errors.New("insufficient allocated quantity")
	}

	// UPDATE QUANTITIES
	newAvailable := item.AvailableQty + quantityToDeallocate
	newAllocated := item.AllocatedQty - quantityToDeallocate

	err = s.repo.UpdateQuantities(id, newAvailable, newAllocated, item.DistributedQty)
	if err != nil {
		log.Printf("INVENTORY: Failed to deallocate inventory: %v", err)
		return err
	}

	log.Printf("INVENTORY: Deallocated %d units from inventory item %d", quantityToDeallocate, id)
	return nil
}

// DELETES INVENTORY ITEM
func (s *Service) Delete(id uint, orgID uint) error {
	item, err := s.GetByID(id, orgID)
	if err != nil {
		return err
	}

	// CHECK IF ANY QUANTITY IS ALLOCATED OR DISTRIBUTED
	if item.AllocatedQty > 0 || item.DistributedQty > 0 {
		return errors.New("cannot delete inventory with allocated or distributed items")
	}

	err = s.repo.Delete(id)
	if err != nil {
		log.Printf("INVENTORY: Failed to delete inventory item: %v", err)
		return err
	}

	log.Printf("INVENTORY: Deleted inventory item %d", id)
	return nil
}

// GETS INVENTORY STATISTICS
func (s *Service) GetStats(orgID uint) (*InventoryStats, error) {
	stats, err := s.repo.GetStats(orgID)
	if err != nil {
		log.Printf("INVENTORY: Failed to get inventory stats: %v", err)
		return nil, err
	}

	log.Printf("INVENTORY: Retrieved inventory stats for org %d", orgID)
	return stats, nil
}
