package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateEstate implements the EstateService.CreateEstate method
func (s *service) CreateEstate(ctx context.Context, width, length int) (uuid.UUID, error) {
	// Validate inputs (although this should also be validated at the API level and DB constraint level)
	if width < 1 || width > 50000 || length < 1 || length > 50000 {
		return uuid.Nil, errors.New("invalid estate dimensions")
	}

	return s.repo.CreateEstate(ctx, width, length)
}

// GetEstate implements the EstateService.GetEstate method
func (s *service) GetEstate(ctx context.Context, id uuid.UUID) (width, length int, err error) {
	// Check if estate exists
	width, length, err = s.repo.GetEstate(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, errors.New("estate not found")
		}
		return 0, 0, err
	}

	return width, length, nil
} 