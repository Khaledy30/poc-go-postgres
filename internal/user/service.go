package user

import (
	"context"
	"errors"
	"poc-go/internal/validator/email"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateUserInput) (*User, error) {
	// Validação de regra de negócio
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("nome é obrigatório")
	}
	if err := email.ValidateEmail(input.Email); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, input)
}

func (s *Service) List(ctx context.Context, page, limit int) ([]User, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.repo.List(ctx, page, limit)
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
