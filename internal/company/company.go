//go:generate mockgen -destination=company_mockgen_test.go -package=company github.com/by-sabbir/company-microservice-rest/internal/company Store

package company

import (
	"context"
	"errors"
	"log"
)

var (
	ErrTypeNotFound = errors.New("provided company type is not defined")
)

var CompanyType = []string{
	"Corporations", "NonProfit", "Cooperative", "Sole Proprietorship",
}

// Company - representation of a company structure
type Company struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	TotalEmployees int    `json:"total_employees,omitempty"`
	IsRegistered   bool   `json:"is_registered,omitempty"`
	Type           string `json:"type,omitempty"`
}

// Service - is the struct containing business logics
type Service struct {
	Store Store
}

// NewService - returns to a pointer to a new service
func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

// Store - this interface defines all the methods to operate
type Store interface {
	GetCompany(context.Context, string) (Company, error)
	PostCompany(context.Context, Company) (Company, error)
	DeleteCompany(context.Context, string) error
	PartialUpdateCompany(context.Context, string, Company) (Company, error)
}

// ScanType - is the replacement for enum in database
func (c *Company) ScanType() error {
	for _, v := range CompanyType {
		if c.Type == v {
			return nil
		}
	}
	return ErrTypeNotFound
}

func (s *Service) GetCompany(ctx context.Context, id string) (Company, error) {
	cmp, err := s.Store.GetCompany(ctx, id)
	if err != nil {
		return Company{}, err
	}
	return cmp, nil
}

func (s *Service) PostCompany(ctx context.Context, cmp Company) (Company, error) {
	cmp, err := s.Store.PostCompany(ctx, cmp)
	if err != nil {
		return Company{}, err
	}
	return cmp, nil
}

func (s *Service) DeleteCompany(ctx context.Context, id string) error {
	err := s.Store.DeleteCompany(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) PartialUpdateCompany(
	ctx context.Context, id string, cmp Company,
) (Company, error) {
	log.Println("Updating Company...")
	cmp, err := s.Store.PartialUpdateCompany(ctx, id, cmp)
	if err != nil {
		return Company{}, err
	}
	return cmp, nil
}
