package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"

	"drone/internal/repository/mocks"
)

func TestCreateEstate(t *testing.T) {
	// Initialize Echo
	e := echo.New()
	
	// Setup test request
	reqBody := `{"width": 10, "length": 20}`
	req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Setup mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	// Create mock repository
	mockRepo := mocks.NewMockRepository(ctrl)
	
	// Create handler with mock repository
	h := NewHandler(mockRepo)
	
	// Perform the test
	err := h.CreateEstate(c)
	
	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestCreateTree(t *testing.T) {
	// Initialize Echo
	e := echo.New()
	
	// Generate estate ID
	googleUUID := uuid.New()
	estateID := openapi_types.UUID(googleUUID)
	
	// Setup test request
	reqBody := `{"x": 3, "y": 4, "height": 5}`
	req := httptest.NewRequest(http.MethodPost, "/estate/"+googleUUID.String()+"/tree", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(googleUUID.String())
	
	// Setup mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	// Create mock repository
	mockRepo := mocks.NewMockRepository(ctrl)
	
	// Create handler with mock repository
	h := NewHandler(mockRepo)
	
	// Perform the test
	err := h.CreateTree(c, estateID)
	
	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestPing(t *testing.T) {
	// Initialize Echo
	e := echo.New()
	
	// Setup test request
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Setup mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	// Create mock repository
	mockRepo := mocks.NewMockRepository(ctrl)
	
	// Create handler with mock repository
	h := NewHandler(mockRepo)
	
	// Perform the test
	err := h.Ping(c)
	
	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "ok")
	assert.Contains(t, rec.Body.String(), "API is running")
} 