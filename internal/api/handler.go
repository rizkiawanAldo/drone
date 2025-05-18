package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"drone/generated"
	"drone/internal/service"
)

// Handler implements the generated ServerInterface
type Handler struct {
	service service.Service
}

// NewHandler creates a new API handler with the given service
func NewHandler(svc service.Service) *Handler {
	return &Handler{
		service: svc,
	}
}

// CreateEstate creates a new estate
func (h *Handler) CreateEstate(ctx echo.Context) error {
	var req generated.EstateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: strPtr("Invalid request format"),
		})
	}

	estateID, err := h.service.CreateEstate(ctx.Request().Context(), int(req.Width), int(req.Length))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: strPtr(err.Error()),
		})
	}

	id := openapi_types.UUID(estateID)
	return ctx.JSON(http.StatusCreated, generated.EstateResponse{
		Id: &id,
	})
}

// CreateTree adds a tree to an estate
func (h *Handler) CreateTree(ctx echo.Context, id openapi_types.UUID) error {
	var req generated.TreeRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: strPtr("Invalid request format"),
		})
	}

	// Since openapi_types.UUID is an alias for uuid.UUID, we can use it directly
	estateID := uuid.UUID(id)

	treeID, err := h.service.CreateTree(ctx.Request().Context(), estateID, int(req.X), int(req.Y), int(req.Height))
	if err != nil {
		if err.Error() == "estate not found" {
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{
				Message: strPtr("Estate not found"),
			})
		}
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: strPtr(err.Error()),
		})
	}

	// Convert uuid.UUID to openapi_types.UUID (just a type conversion since they're the same underlying type)
	treeUUID := openapi_types.UUID(treeID)
	return ctx.JSON(http.StatusCreated, generated.TreeResponse{
		Id: &treeUUID,
	})
}

// GetEstateStats gets stats about trees in an estate
func (h *Handler) GetEstateStats(ctx echo.Context, id openapi_types.UUID) error {
	// Since openapi_types.UUID is an alias for uuid.UUID, we can use it directly
	estateID := uuid.UUID(id)

	count, maxHeight, minHeight, medianHeight, err := h.service.GetTreeStats(ctx.Request().Context(), estateID)
	if err != nil {
		if err.Error() == "estate not found" {
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{
				Message: strPtr("Estate not found"),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: strPtr(err.Error()),
		})
	}

	count32 := int32(count)
	maxHeight32 := int32(maxHeight)
	minHeight32 := int32(minHeight)
	medianHeight32 := int32(medianHeight)

	return ctx.JSON(http.StatusOK, generated.StatsResponse{
		Count:        &count32,
		MaxHeight:    &maxHeight32,
		MinHeight:    &minHeight32,
		MedianHeight: &medianHeight32,
	})
}

// GetDronePlan gets the drone monitoring travel plan
func (h *Handler) GetDronePlan(ctx echo.Context, id openapi_types.UUID, params generated.GetDronePlanParams) error {
	// Since openapi_types.UUID is an alias for uuid.UUID, we can use it directly
	estateID := uuid.UUID(id)

	var response generated.DronePlanResponse

	// If max_distance is provided, calculate with rest
	if params.MaxDistance != nil {
		maxDistance := int(*params.MaxDistance)
		
		if maxDistance <= 0 {
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
				Message: strPtr("Max distance must be positive"),
			})
		}

		distance, restX, restY, err := h.service.CalculateDronePathWithRest(ctx.Request().Context(), estateID, maxDistance)
		if err != nil {
			if err.Error() == "estate not found" {
				return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{
					Message: strPtr("Estate not found"),
				})
			}
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
				Message: strPtr(err.Error()),
			})
		}

		// Convert to int32 for the response
		distance32 := int32(distance)
		restX32 := int32(restX)
		restY32 := int32(restY)

		response = generated.DronePlanResponse{
			Distance: &distance32,
			Rest: &struct {
				X *int32 `json:"x,omitempty"`
				Y *int32 `json:"y,omitempty"`
			}{
				X: &restX32,
				Y: &restY32,
			},
		}
	} else {
		// Calculate without rest
		distance, err := h.service.CalculateDronePath(ctx.Request().Context(), estateID)
		if err != nil {
			if err.Error() == "estate not found" {
				return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{
					Message: strPtr("Estate not found"),
				})
			}
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
				Message: strPtr(err.Error()),
			})
		}

		// Convert to int32 for the response
		distance32 := int32(distance)
		response = generated.DronePlanResponse{
			Distance: &distance32,
		}
	}

	return ctx.JSON(http.StatusOK, response)
}

// ListEstates lists all estates
func (h *Handler) ListEstates(ctx echo.Context) error {
	estates, err := h.service.ListEstates(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: strPtr(err.Error()),
		})
	}

	// Convert from repository.Estate to generated.EstateListItem
	response := make([]generated.EstateListItem, len(estates))
	for i, estate := range estates {
		id := openapi_types.UUID(estate.ID)
		width := int32(estate.Width)
		length := int32(estate.Length)
		
		response[i] = generated.EstateListItem{
			Id:     &id,
			Width:  &width,
			Length: &length,
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

// Helper function to convert a string to a pointer
func strPtr(s string) *string {
	return &s
} 