package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"drone/generated"
	"drone/internal/repository"
)

// Handler implements the generated ServerInterface
type Handler struct {
	repo repository.Repository
}

// NewHandler creates a new API handler with the given repository
func NewHandler(repo repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// CreateEstate creates a new estate
func (h *Handler) CreateEstate(ctx echo.Context) error {
	// Empty stub implementation
	// We'll implement the business logic later
	return ctx.JSON(http.StatusCreated, generated.EstateResponse{
		Id: uuid.New().String(),
	})
}

// CreateTree adds a tree to an estate
func (h *Handler) CreateTree(ctx echo.Context, id string) error {
	// Empty stub implementation
	// We'll implement the business logic later
	return ctx.JSON(http.StatusCreated, generated.TreeResponse{
		Id: uuid.New().String(),
	})
}

// GetEstateStats gets stats about trees in an estate
func (h *Handler) GetEstateStats(ctx echo.Context, id string) error {
	// Empty stub implementation
	// We'll implement the business logic later
	return ctx.JSON(http.StatusOK, generated.StatsResponse{
		Count:        0,
		MaxHeight:    0,
		MinHeight:    0,
		MedianHeight: 0,
	})
}

// GetDronePlan gets the drone monitoring travel plan
func (h *Handler) GetDronePlan(ctx echo.Context, id string, params generated.GetDronePlanParams) error {
	// Empty stub implementation
	// We'll implement the business logic later
	response := generated.DronePlanResponse{
		Distance: 0,
	}
	
	// If max_distance is provided, include rest info
	if params.MaxDistance != nil {
		response.Rest = &generated.DronePlanResponseRest{
			X: 1,
			Y: 1,
		}
	}
	
	return ctx.JSON(http.StatusOK, response)
} 