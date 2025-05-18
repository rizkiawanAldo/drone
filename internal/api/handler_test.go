package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"

	"drone/generated"
	"drone/internal/repository"
	"drone/internal/repository/mocks"
)

func TestCreateEstate(t *testing.T) {
	testCases := []struct {
		name           string
		requestBody    string
		mockSetup      func(*mocks.MockRepository)
		expectedStatus int
		checkResponse  func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:        "Success",
			requestBody: `{"width": 1000, "length": 2000}`,
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().
					CreateEstate(gomock.Any(), 1000, 2000).
					Return(uuid.New(), nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response generated.EstateResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response.Id)
			},
		},
		{
			name:        "Invalid Request - Missing Fields",
			requestBody: `{"width": 1000}`,
			mockSetup:   func(mockRepo *mocks.MockRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Invalid Request - Invalid Width",
			requestBody: `{"width": 0, "length": 2000}`,
			mockSetup:   func(mockRepo *mocks.MockRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Invalid Request - Invalid Length",
			requestBody: `{"width": 1000, "length": 60000}`,
			mockSetup:   func(mockRepo *mocks.MockRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Repository Error",
			requestBody: `{"width": 1000, "length": 2000}`,
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().
					CreateEstate(gomock.Any(), 1000, 2000).
					Return(uuid.UUID{}, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize Echo
			e := echo.New()
			
			// Setup test request
			req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			
			// Create mock repository
			mockRepo := mocks.NewMockRepository(ctrl)
			
			// Setup mock expectations
			tc.mockSetup(mockRepo)
			
			// Create handler with mock repository
			h := NewHandler(mockRepo)
			
			// Perform the test
			_ = h.CreateEstate(c)
			
			// Assert the results
			assert.Equal(t, tc.expectedStatus, rec.Code)
			
			// Additional response checks if provided
			if tc.checkResponse != nil {
				tc.checkResponse(t, rec)
			}
		})
	}
}

func TestCreateTree(t *testing.T) {
	// Generate estate ID
	estateID := uuid.New()
	estateUUID := openapi_types.UUID(estateID)

	testCases := []struct {
		name           string
		requestBody    string
		mockSetup      func(*mocks.MockRepository)
		expectedStatus int
		checkResponse  func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:        "Success",
			requestBody: `{"x": 5, "y": 10, "height": 15}`,
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().
					GetEstate(gomock.Any(), estateID).
					Return(100, 100, nil)
				mockRepo.EXPECT().
					CreateTree(gomock.Any(), estateID, 5, 10, 15).
					Return(uuid.New(), nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response generated.TreeResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response.Id)
			},
		},
		{
			name:        "Invalid Request - Missing Fields",
			requestBody: `{"x": 5, "y": 10}`,
			mockSetup:   func(mockRepo *mocks.MockRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Invalid Request - Invalid X",
			requestBody: `{"x": 0, "y": 10, "height": 15}`,
			mockSetup:   func(mockRepo *mocks.MockRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Invalid Request - Invalid Y",
			requestBody: `{"x": 5, "y": 0, "height": 15}`,
			mockSetup:   func(mockRepo *mocks.MockRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Invalid Request - Invalid Height",
			requestBody: `{"x": 5, "y": 10, "height": 40}`,
			mockSetup:   func(mockRepo *mocks.MockRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Estate Not Found",
			requestBody: `{"x": 5, "y": 10, "height": 15}`,
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().
					GetEstate(gomock.Any(), estateID).
					Return(0, 0, errors.New("estate not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:        "Tree Out of Bounds",
			requestBody: `{"x": 200, "y": 10, "height": 15}`,
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().
					GetEstate(gomock.Any(), estateID).
					Return(100, 100, nil)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Repository Error",
			requestBody: `{"x": 5, "y": 10, "height": 15}`,
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().
					GetEstate(gomock.Any(), estateID).
					Return(100, 100, nil)
				mockRepo.EXPECT().
					CreateTree(gomock.Any(), estateID, 5, 10, 15).
					Return(uuid.UUID{}, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize Echo
			e := echo.New()
			
			// Setup test request
			req := httptest.NewRequest(http.MethodPost, "/estate/"+estateID.String()+"/tree", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(estateID.String())
			
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			
			// Create mock repository
			mockRepo := mocks.NewMockRepository(ctrl)
			
			// Setup mock expectations
			tc.mockSetup(mockRepo)
			
			// Create handler with mock repository
			h := NewHandler(mockRepo)
			
			// Perform the test
			_ = h.CreateTree(c, estateUUID)
			
			// Assert the results
			assert.Equal(t, tc.expectedStatus, rec.Code)
			
			// Additional response checks if provided
			if tc.checkResponse != nil {
				tc.checkResponse(t, rec)
			}
		})
	}
}

func TestGetEstateStats(t *testing.T) {
	estateID := uuid.New()
	estateUUID := openapi_types.UUID(estateID)

	trees := []repository.Tree{
		{ID: uuid.New(), EstateID: estateID, X: 5, Y: 10, Height: 10},
		{ID: uuid.New(), EstateID: estateID, X: 15, Y: 20, Height: 15},
		{ID: uuid.New(), EstateID: estateID, X: 25, Y: 30, Height: 20},
	}

	testCases := []struct {
		name           string
		mockSetup      func(*mocks.MockRepository)
		expectedStatus int
		checkResponse  func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().
					GetEstate(gomock.Any(), estateID).
					Return(100, 100, nil)
				mockRepo.EXPECT().
					GetTrees(gomock.Any(), estateID).
					Return(trees, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response generated.StatsResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, int32(3), *response.Count)
				assert.Equal(t, int32(20), *response.MaxHeight)
				assert.Equal(t, int32(10), *response.MinHeight)
				assert.Equal(t, int32(15), *response.MedianHeight)
			},
		},
		{
			name: "Success - No Trees",
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().
					GetEstate(gomock.Any(), estateID).
					Return(100, 100, nil)
				mockRepo.EXPECT().
					GetTrees(gomock.Any(), estateID).
					Return([]repository.Tree{}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response generated.StatsResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, int32(0), *response.Count)
				assert.Equal(t, int32(0), *response.MaxHeight)
				assert.Equal(t, int32(0), *response.MinHeight)
				assert.Equal(t, int32(0), *response.MedianHeight)
			},
		},
		{
			name: "Estate Not Found",
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().
					GetEstate(gomock.Any(), estateID).
					Return(0, 0, errors.New("estate not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Repository Error",
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().
					GetEstate(gomock.Any(), estateID).
					Return(100, 100, nil)
				mockRepo.EXPECT().
					GetTrees(gomock.Any(), estateID).
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize Echo
			e := echo.New()
			
			// Setup test request
			req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID.String()+"/stats", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(estateID.String())
			
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			
			// Create mock repository
			mockRepo := mocks.NewMockRepository(ctrl)
			
			// Setup mock expectations
			tc.mockSetup(mockRepo)
			
			// Create handler with mock repository
			h := NewHandler(mockRepo)
			
			// Perform the test
			_ = h.GetEstateStats(c, estateUUID)
			
			// Assert the results
			assert.Equal(t, tc.expectedStatus, rec.Code)
			
			// Additional response checks if provided
			if tc.checkResponse != nil {
				tc.checkResponse(t, rec)
			}
		})
	}
}

func TestGetDronePlan(t *testing.T) {
	estateID := uuid.New()
	estateUUID := openapi_types.UUID(estateID)

	trees := []repository.Tree{
		{ID: uuid.New(), EstateID: estateID, X: 5, Y: 10, Height: 10},
		{ID: uuid.New(), EstateID: estateID, X: 15, Y: 20, Height: 15},
		{ID: uuid.New(), EstateID: estateID, X: 25, Y: 30, Height: 20},
	}

	testCases := []struct {
		name           string
		maxDistance    *int32
		mockSetup      func(*mocks.MockRepository)
		expectedStatus int
		checkResponse  func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "Success - No Max Distance",
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().
					GetEstate(gomock.Any(), estateID).
					Return(100, 100, nil)
				mockRepo.EXPECT().
					GetTrees(gomock.Any(), estateID).
					Return(trees, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response generated.DronePlanResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response.Distance)
				assert.Nil(t, response.Rest)
			},
		},
		{
			name: "Success - With Max Distance",
			maxDistance: func() *int32 { val := int32(50); return &val }(),
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().
					GetEstate(gomock.Any(), estateID).
					Return(100, 100, nil)
				mockRepo.EXPECT().
					GetTrees(gomock.Any(), estateID).
					Return(trees, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response generated.DronePlanResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response.Distance)
				assert.NotNil(t, response.Rest)
				assert.NotNil(t, response.Rest.X)
				assert.NotNil(t, response.Rest.Y)
			},
		},
		{
			name: "Invalid Max Distance",
			maxDistance: func() *int32 { val := int32(0); return &val }(),
			mockSetup: func(mockRepo *mocks.MockRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Estate Not Found",
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().
					GetEstate(gomock.Any(), estateID).
					Return(0, 0, errors.New("estate not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Repository Error",
			mockSetup: func(mockRepo *mocks.MockRepository) {
				mockRepo.EXPECT().
					GetEstate(gomock.Any(), estateID).
					Return(100, 100, nil)
				mockRepo.EXPECT().
					GetTrees(gomock.Any(), estateID).
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize Echo
			e := echo.New()
			
			// Setup test request
			path := "/estate/" + estateID.String() + "/drone-plan"
			if tc.maxDistance != nil {
				path += "?max_distance=" + string(*tc.maxDistance)
			}
			
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(estateID.String())
			
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			
			// Create mock repository
			mockRepo := mocks.NewMockRepository(ctrl)
			
			// Setup mock expectations
			tc.mockSetup(mockRepo)
			
			// Create handler with mock repository
			h := NewHandler(mockRepo)
			
			// Prepare params
			params := generated.GetDronePlanParams{
				MaxDistance: tc.maxDistance,
			}
			
			// Perform the test
			_ = h.GetDronePlan(c, estateUUID, params)
			
			// Assert the results
			assert.Equal(t, tc.expectedStatus, rec.Code)
			
			// Additional response checks if provided
			if tc.checkResponse != nil {
				tc.checkResponse(t, rec)
			}
		})
	}
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