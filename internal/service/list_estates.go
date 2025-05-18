package service

import (
	"context"

	"drone/internal/repository"
)

// ListEstates implements the EstateService.ListEstates method
func (s *service) ListEstates(ctx context.Context) ([]repository.Estate, error) {
	return s.repo.ListEstates(ctx)
} 