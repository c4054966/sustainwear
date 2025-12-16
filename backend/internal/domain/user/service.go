package user

import (
	"fmt"
	"sustainwear/pkg/validator"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CREATES NEW USER ACCOUNT
func (s *Service) Register(req *RegisterRequest) (*User, error) {
	if err := validator.ValidateRequired(map[string]string{
		"email":     req.Email,
		"password":  req.Password,
		"full_name": req.FullName,
		"role":      req.Role,
	}); err != nil {
		return nil, err
	}

	if !validator.IsValidEmail(req.Email) {
		return nil, fmt.Errorf("invalid email format")
	}

	if valid, err := validator.IsValidPassword(req.Password); !valid {
		return nil, fmt.Errorf("%v", err)
	}

	if !validator.IsValidRole(req.Role) {
		return nil, fmt.Errorf("invalid role: must be donor, charity_staff, or admin")
	}

	existing, _ := s.repo.GetByEmail(req.Email)
	if existing != nil {
		return nil, fmt.Errorf("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &User{
		Email:          req.Email,
		PasswordHash:   string(hashedPassword),
		FullName:       req.FullName,
		Role:           req.Role,
		OrganisationID: req.OrganisationID,
		IsActive:       true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("server error: %w", err)
	}

	return user, nil
}

// AUTHENTICATES USER AND RETURNS USER DETAILS
func (s *Service) Login(email, password string) (*User, error) {
	if validator.IsEmpty(email) || validator.IsEmpty(password) {
		return nil, fmt.Errorf("email and password are required")
	}

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// FETCH USER DETAILS BY ID
func (s *Service) GetByID(id uint) (*User, error) {
	return s.repo.GetByID(id)
}

// FETCH USER DETAILS BY EMAIL
func (s *Service) GetByEmail(email string) (*User, error) {
	return s.repo.GetByEmail(email)
}

// LIST ALL USERS
func (s *Service) List() ([]User, error) {
	return s.repo.List()
}

func (s *Service) ListPaginated(limit, offset int) ([]User, error) {
	return s.repo.ListPaginated(limit, offset)
}

// UPDATE USER DETAILS
func (s *Service) Update(user *User) error {
	if !validator.IsValidRole(user.Role) {
		return fmt.Errorf("invalid role")
	}

	return s.repo.Update(user)
}

// DELETE USER BY ID
func (s *Service) Delete(id uint) error {
	return s.repo.Delete(id)
}
