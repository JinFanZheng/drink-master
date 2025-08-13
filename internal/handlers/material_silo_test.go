package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/enums"
	"github.com/ddteam/drink-master/internal/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = models.AutoMigrate(db)
	require.NoError(t, err)

	return db
}

func TestMaterialSiloHandler_GetPaging(t *testing.T) {
	db := setupTestDB(t)
	handler := NewMaterialSiloHandler(db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/MaterialSilo/GetPaging", handler.GetPaging)

	// Create test machine
	machine := &models.Machine{
		ID:             "test_machine",
		MachineOwnerId: "test_owner",
		MachineNo:      "M001",
		Name:           "Test Machine",
		BusinessStatus: enums.BusinessStatusOpen,
	}
	require.NoError(t, db.Create(machine).Error)

	// Create test product
	product := &models.Product{
		ID:       "test_product",
		Name:     "Test Product",
		Category: func() *string { s := "drink"; return &s }(),
	}
	require.NoError(t, db.Create(product).Error)

	// Create test material silos
	silos := []*models.MaterialSilo{
		{
			ID:          "silo_1",
			MachineID:   machine.ID,
			SiloNo:      1,
			ProductID:   &product.ID,
			Stock:       50,
			MaxCapacity: 100,
			SaleStatus:  enums.SaleStatusOn,
		},
		{
			ID:          "silo_2",
			MachineID:   machine.ID,
			SiloNo:      2,
			Stock:       30,
			MaxCapacity: 100,
			SaleStatus:  enums.SaleStatusOff,
		},
	}

	for _, silo := range silos {
		require.NoError(t, db.Create(silo).Error)
	}

	tests := []struct {
		name           string
		request        contracts.GetMaterialSiloPagingRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful paging request",
			request: contracts.GetMaterialSiloPagingRequest{
				MachineID: machine.ID,
				PageIndex: 1,
				PageSize:  10,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "invalid request - missing machine id",
			request: contracts.GetMaterialSiloPagingRequest{
				PageIndex: 1,
				PageSize:  10,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "invalid request - zero page index",
			request: contracts.GetMaterialSiloPagingRequest{
				MachineID: machine.ID,
				PageIndex: 0,
				PageSize:  10,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "machine not found",
			request: contracts.GetMaterialSiloPagingRequest{
				MachineID: "nonexistent",
				PageIndex: 1,
				PageSize:  10,
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare request body
			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest(http.MethodPost, "/api/MaterialSilo/GetPaging", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check response status
			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectError && w.Code == http.StatusOK {
				// Parse successful response
				var response contracts.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.True(t, response.Success)
				assert.NotNil(t, response.Data)

				// Convert data to MaterialSiloPaging
				dataBytes, _ := json.Marshal(response.Data)
				var paging contracts.MaterialSiloPaging
				err = json.Unmarshal(dataBytes, &paging)
				require.NoError(t, err)

				assert.Len(t, paging.Items, 2)
				assert.Equal(t, int64(2), paging.TotalCount)
				assert.Equal(t, 1, paging.PageIndex)
				assert.Equal(t, 10, paging.PageSize)
			}
		})
	}
}

func TestMaterialSiloHandler_UpdateStock(t *testing.T) {
	db := setupTestDB(t)
	handler := NewMaterialSiloHandler(db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/MaterialSilo/UpdateStock", handler.UpdateStock)

	// Create test material silo
	silo := &models.MaterialSilo{
		ID:          "test_silo",
		MachineID:   "test_machine",
		SiloNo:      1,
		Stock:       50,
		MaxCapacity: 100,
		SaleStatus:  enums.SaleStatusOn,
	}
	require.NoError(t, db.Create(silo).Error)

	tests := []struct {
		name           string
		request        contracts.UpdateMaterialSiloStockRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful stock update",
			request: contracts.UpdateMaterialSiloStockRequest{
				ID:    "test_silo",
				Stock: 75,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "invalid request - missing id",
			request: contracts.UpdateMaterialSiloStockRequest{
				Stock: 75,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "stock exceeds capacity",
			request: contracts.UpdateMaterialSiloStockRequest{
				ID:    "test_silo",
				Stock: 150,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "silo not found",
			request: contracts.UpdateMaterialSiloStockRequest{
				ID:    "nonexistent",
				Stock: 50,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare request body
			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest(http.MethodPost, "/api/MaterialSilo/UpdateStock", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check response status
			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectError && w.Code == http.StatusOK {
				// Parse successful response
				var response contracts.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.True(t, response.Success)
			}
		})
	}
}

func TestMaterialSiloHandler_UpdateProduct(t *testing.T) {
	db := setupTestDB(t)
	handler := NewMaterialSiloHandler(db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/MaterialSilo/UpdateProduct", handler.UpdateProduct)

	// Create test product
	product := &models.Product{
		ID:       "test_product",
		Name:     "Test Product",
		Category: func() *string { s := "drink"; return &s }(),
	}
	require.NoError(t, db.Create(product).Error)

	// Create test material silo
	silo := &models.MaterialSilo{
		ID:          "test_silo",
		MachineID:   "test_machine",
		SiloNo:      1,
		Stock:       50,
		MaxCapacity: 100,
		SaleStatus:  enums.SaleStatusOn,
	}
	require.NoError(t, db.Create(silo).Error)

	tests := []struct {
		name           string
		request        contracts.UpdateMaterialSiloProductRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful product update",
			request: contracts.UpdateMaterialSiloProductRequest{
				ID:        "test_silo",
				ProductID: "test_product",
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "invalid request - missing id",
			request: contracts.UpdateMaterialSiloProductRequest{
				ProductID: "test_product",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "invalid request - missing product id",
			request: contracts.UpdateMaterialSiloProductRequest{
				ID: "test_silo",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "silo not found",
			request: contracts.UpdateMaterialSiloProductRequest{
				ID:        "nonexistent",
				ProductID: "test_product",
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name: "product not found",
			request: contracts.UpdateMaterialSiloProductRequest{
				ID:        "test_silo",
				ProductID: "nonexistent",
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare request body
			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest(http.MethodPost, "/api/MaterialSilo/UpdateProduct", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check response status
			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectError && w.Code == http.StatusOK {
				// Parse successful response
				var response contracts.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.True(t, response.Success)
			}
		})
	}
}

func TestMaterialSiloHandler_ToggleSaleStatus(t *testing.T) {
	db := setupTestDB(t)
	handler := NewMaterialSiloHandler(db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/MaterialSilo/ToggleSaleStatus", handler.ToggleSaleStatus)

	// Create test product
	product := &models.Product{
		ID:       "test_product",
		Name:     "Test Product",
		Category: func() *string { s := "drink"; return &s }(),
	}
	require.NoError(t, db.Create(product).Error)

	// Create test material silos for different scenarios
	silos := []*models.MaterialSilo{
		{
			ID:          "silo_with_product_and_stock",
			MachineID:   "test_machine",
			SiloNo:      1,
			ProductID:   &product.ID,
			Stock:       50,
			MaxCapacity: 100,
			SaleStatus:  enums.SaleStatusOff,
		},
		{
			ID:          "silo_without_product",
			MachineID:   "test_machine",
			SiloNo:      2,
			ProductID:   nil,
			Stock:       50,
			MaxCapacity: 100,
			SaleStatus:  enums.SaleStatusOff,
		},
		{
			ID:          "silo_without_stock",
			MachineID:   "test_machine",
			SiloNo:      3,
			ProductID:   &product.ID,
			Stock:       0,
			MaxCapacity: 100,
			SaleStatus:  enums.SaleStatusOff,
		},
	}

	for _, silo := range silos {
		require.NoError(t, db.Create(silo).Error)
	}

	tests := []struct {
		name           string
		request        contracts.ToggleSaleMaterialSiloRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful turn on sale status",
			request: contracts.ToggleSaleMaterialSiloRequest{
				ID:         "silo_with_product_and_stock",
				SaleStatus: contracts.SaleStatusOn,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "successful turn off sale status",
			request: contracts.ToggleSaleMaterialSiloRequest{
				ID:         "silo_with_product_and_stock",
				SaleStatus: contracts.SaleStatusOff,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "turn on without product",
			request: contracts.ToggleSaleMaterialSiloRequest{
				ID:         "silo_without_product",
				SaleStatus: contracts.SaleStatusOn,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "turn on without stock",
			request: contracts.ToggleSaleMaterialSiloRequest{
				ID:         "silo_without_stock",
				SaleStatus: contracts.SaleStatusOn,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "invalid request - missing id",
			request: contracts.ToggleSaleMaterialSiloRequest{
				SaleStatus: contracts.SaleStatusOn,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "invalid request - invalid sale status",
			request: contracts.ToggleSaleMaterialSiloRequest{
				ID:         "silo_with_product_and_stock",
				SaleStatus: "Invalid",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "silo not found",
			request: contracts.ToggleSaleMaterialSiloRequest{
				ID:         "nonexistent",
				SaleStatus: contracts.SaleStatusOn,
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare request body
			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest(http.MethodPost, "/api/MaterialSilo/ToggleSaleStatus", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check response status
			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectError && w.Code == http.StatusOK {
				// Parse successful response
				var response contracts.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.True(t, response.Success)
			}
		})
	}
}
