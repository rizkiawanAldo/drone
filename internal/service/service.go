package service

import (
	"context"

	"github.com/google/uuid"

	"drone/internal/repository"
)

// EstateService defines the interface for estate-related operations
type EstateService interface {
	CreateEstate(ctx context.Context, width, length int) (uuid.UUID, error)
	GetEstate(ctx context.Context, id uuid.UUID) (width, length int, err error)
	ListEstates(ctx context.Context) ([]repository.Estate, error)
}

// TreeService defines the interface for tree-related operations
type TreeService interface {
	CreateTree(ctx context.Context, estateID uuid.UUID, x, y, height int) (uuid.UUID, error)
	GetTreeStats(ctx context.Context, estateID uuid.UUID) (count, maxHeight, minHeight, medianHeight int, err error)
}

// DroneService defines the interface for drone-related operations
type DroneService interface {
	CalculateDronePath(ctx context.Context, estateID uuid.UUID) (distance int, err error)
	CalculateDronePathWithRest(ctx context.Context, estateID uuid.UUID, maxDistance int) (distance int, restX, restY int, err error)
}

// Service combines all service interfaces
type Service interface {
	EstateService
	TreeService
	DroneService
}

// service implements the Service interface
type service struct {
	repo repository.Repository
}

// NewService creates a new service with the given repository
func NewService(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
} 