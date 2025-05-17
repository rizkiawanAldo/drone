package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

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
	googleUUID := uuid.New()
	id := openapi_types.UUID(googleUUID)
	return ctx.JSON(http.StatusCreated, generated.EstateResponse{
		Id: &id,
	})
}

// CreateTree adds a tree to an estate
func (h *Handler) CreateTree(ctx echo.Context, id openapi_types.UUID) error {
	// Empty stub implementation
	// We'll implement the business logic later
	googleUUID := uuid.New()
	treeId := openapi_types.UUID(googleUUID)
	return ctx.JSON(http.StatusCreated, generated.TreeResponse{
		Id: &treeId,
	})
}

// GetEstateStats gets stats about trees in an estate
func (h *Handler) GetEstateStats(ctx echo.Context, id openapi_types.UUID) error {
	// Empty stub implementation
	// We'll implement the business logic later
	count := int32(0)
	maxHeight := int32(0)
	minHeight := int32(0)
	medianHeight := int32(0)
	return ctx.JSON(http.StatusOK, generated.StatsResponse{
		Count:        &count,
		MaxHeight:    &maxHeight,
		MinHeight:    &minHeight,
		MedianHeight: &medianHeight,
	})
}

// GetDronePlan gets the drone monitoring travel plan
func (h *Handler) GetDronePlan(ctx echo.Context, id openapi_types.UUID, params generated.GetDronePlanParams) error {
	// Empty stub implementation
	// We'll implement the business logic later
	distance := int32(0)
	response := generated.DronePlanResponse{
		Distance: &distance,
	}
	
	// If max_distance is provided, include rest info
	if params.MaxDistance != nil {
		x := int32(1)
		y := int32(1)
		response.Rest = &struct {
			X *int32 `json:"x,omitempty"`
			Y *int32 `json:"y,omitempty"`
		}{
			X: &x,
			Y: &y,
		}
	}
	
	return ctx.JSON(http.StatusOK, response)
}

// Ping handles ping requests to check if the API is available
func (h *Handler) Ping(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"message": "API is running",
	})
} 