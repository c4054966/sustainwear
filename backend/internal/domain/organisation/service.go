package organisation

import (
	"errors"
	"log"
	"regexp"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CREATES ORGANISATION
func (s *Service) Create(req *CreateOrgRequest) (*Organisation, error) {
	if req.Name == "" {
		return nil, errors.New("organisation name is required")
	}

	if req.Email == "" {
		return nil, errors.New("email is required")
	}

	if !isValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	if req.Type == "" {
		return nil, errors.New("organisation type is required")
	}

	validTypes := map[string]bool{"charity": true, "ngo": true, "community": true, "religious": true}
	if !validTypes[req.Type] {
		return nil, errors.New("invalid organisation type")
	}

	if req.PostCode != "" && !isValidUKPostcode(req.PostCode) {
		return nil, errors.New("invalid UK postcode format")
	}

	existing, _ := s.repo.GetByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("organisation with this email already exists")
	}

	org := &Organisation{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Email:       req.Email,
		Phone:       req.Phone,
		Address:     req.Address,
		City:        req.City,
		County:      req.County,
		PostCode:    req.PostCode,
		Country:     req.Country,
		Website:     req.Website,
		Status:      "active",
	}

	err := s.repo.Create(org)
	if err != nil {
		log.Printf("ORGANISATION: Failed to create organisation: %v", err)
		return nil, err
	}

	log.Printf("ORGANISATION: Created organisation %s (ID: %d)", org.Name, org.ID)
	return org, nil
}

// GETS ORGANISATION BY ID
func (s *Service) GetByID(id uint) (*Organisation, error) {
	org, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// GETS ORGANISATION BY EMAIL
func (s *Service) GetByEmail(email string) (*Organisation, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	org, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// LISTS ORGANISATIONS
func (s *Service) List(filters map[string]interface{}) ([]Organisation, error) {
	orgs, err := s.repo.List(filters)
	if err != nil {
		log.Printf("ORGANISATION: Failed to list organisations: %v", err)
		return nil, err
	}

	log.Printf("ORGANISATION: Listed %d organisations", len(orgs))
	return orgs, nil
}

// UPDATES ORGANISATION
func (s *Service) Update(id uint, req *UpdateOrgRequest) (*Organisation, error) {
	org, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		if *req.Name == "" {
			return nil, errors.New("organisation name cannot be empty")
		}
		org.Name = *req.Name
	}

	if req.Description != nil {
		org.Description = *req.Description
	}

	if req.Phone != nil {
		org.Phone = *req.Phone
	}

	if req.Address != nil {
		org.Address = *req.Address
	}

	if req.City != nil {
		org.City = *req.City
	}

	if req.County != nil {
		org.County = *req.County
	}

	if req.PostCode != nil {
		if *req.PostCode != "" && !isValidUKPostcode(*req.PostCode) {
			return nil, errors.New("invalid UK postcode format")
		}
		org.PostCode = *req.PostCode
	}

	if req.Website != nil {
		org.Website = *req.Website
	}

	if req.Status != nil {
		validStatuses := map[string]bool{"active": true, "inactive": true, "pending": true}
		if !validStatuses[*req.Status] {
			return nil, errors.New("invalid status")
		}
		org.Status = *req.Status
	}

	err = s.repo.Update(org)
	if err != nil {
		log.Printf("ORGANISATION: Failed to update organisation: %v", err)
		return nil, err
	}

	log.Printf("ORGANISATION: Updated organisation %d", id)
	return org, nil
}

// DELETES ORGANISATION
func (s *Service) Delete(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	err = s.repo.Delete(id)
	if err != nil {
		log.Printf("ORGANISATION: Failed to delete organisation: %v", err)
		return err
	}

	log.Printf("ORGANISATION: Deleted organisation %d", id)
	return nil
}

// GETS ORGANISATION STATS
func (s *Service) GetStats(orgID uint) (*OrgStats, error) {
	_, err := s.repo.GetByID(orgID)
	if err != nil {
		return nil, err
	}

	stats, err := s.repo.GetStats(orgID)
	if err != nil {
		log.Printf("ORGANISATION: Failed to get stats: %v", err)
		return nil, err
	}

	log.Printf("ORGANISATION: Retrieved stats for organisation %d", orgID)
	return stats, nil
}

// VALIDATES EMAIL FORMAT
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// VALIDATES UK POSTCODE FORMAT
func isValidUKPostcode(postcode string) bool {
	postcodeRegex := regexp.MustCompile(`^[A-Z]{1,2}\d{1,2}[A-Z]?\s?\d[A-Z]{2}$`)
	return postcodeRegex.MatchString(postcode)
}
