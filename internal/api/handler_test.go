package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"drone/internal/repository"
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
	estateID := uuid.New().String()
	
	// Setup test request
	reqBody := `{"x": 3, "y": 4, "height": 5}`
	req := httptest.NewRequest(http.MethodPost, "/estate/"+estateID+"/tree", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(estateID)
	
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