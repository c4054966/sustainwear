package donation

import (
	"fmt"
	"sustainwear/pkg/validator"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CREATES NEW DONATION
func (s *Service) Create(donorID uint, req *CreateRequest) (*Donation, error) {
	if err := validator.ValidateRequired(map[string]string{
		"item_name": req.ItemName,
		"category":  req.Category,
		"condition": req.Condition,
	}); err != nil {
		return nil, err
	}

	if req.Quantity < 1 {
		return nil, fmt.Errorf("quantity must be at least 1")
	}

	imagesJSON := MarshalImages(req.Images)

	donation := &Donation{
		DonorID:     donorID,
		ItemName:    req.ItemName,
		Description: req.Description,
		Category:    req.Category,
		Size:        req.Size,
		Gender:      req.Gender,
		Condition:   req.Condition,
		Quantity:    req.Quantity,
		Images:      imagesJSON,
		Status:      StatusPending,
	}

	if err := s.repo.Create(donation); err != nil {
		return nil, fmt.Errorf("server error: %v", err)
	}

	return donation, nil
}

// GETS DONATION BY ID
func (s *Service) GetByID(id uint) (*Donation, error) {
	return s.repo.GetByID(id)
}

// LISTS DONATIONS WITH FILTERS
func (s *Service) List(filters map[string]interface{}) ([]Donation, error) {
	return s.repo.List(filters)
}

// UPDATES DONATION STATUS
func (s *Service) UpdateStatus(id uint, status, notes string) error {
	if !validator.IsValidStatus(status) {
		return fmt.Errorf("invalid status: must be pending, approved, rejected, or distributed")
	}

	return s.repo.UpdateStatus(id, status, notes)
}

// DELETES DONATION BY ID
func (s *Service) Delete(id uint) error {
	return s.repo.Delete(id)
}

// GETS DONATIONS BY DONOR ID
func (s *Service) GetByDonorID(donorID uint) ([]Donation, error) {
	filters := map[string]interface{}{
		"donor_id": donorID,
	}
	return s.repo.List(filters)
}
