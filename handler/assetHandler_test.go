package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"inventory-tracker/models"
	service "inventory-tracker/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func strPtr(s string) *string { return &s } // Go has no literal address-of for a string constant; wrap it in a variable first

// fakeAssetService is a test double for service.AssetService — lets us control
// exactly what AddAsset returns without touching a real or mocked database.
type fakeAssetService struct {
	addAssetFunc func(ctx context.Context, asset models.AddAssetReq) error
	getAssetList func(ctx context.Context, params models.AssetListParams) ([]models.AssetResponse, int, error)
}

func (f *fakeAssetService) AddAsset(ctx context.Context, asset models.AddAssetReq) error {
	return f.addAssetFunc(ctx, asset)
}

func (f *fakeAssetService) GetAssetList(ctx context.Context, params models.AssetListParams) ([]models.AssetResponse, int, error) {
	return f.getAssetList(ctx, params)
}

func TestAddAsset(t *testing.T) {
	gin.SetMode(gin.TestMode) // silences gin logging, makes output readable; doesn't affect behavior

	cases := []struct {
		name           string
		input          models.AddAssetReq
		mockSetup      func() service.AssetService // builds the fake service for this case
		expectedStatus int
	}{
		{
			name: "valid asset",
			input: models.AddAssetReq{
				AssetCode:  "AST-001",
				AssetName:  "Dell Laptop",
				Brand:      strPtr("Dell"),
				CategoryId: 2,
				Status:     "active",
			},
			mockSetup: func() service.AssetService {
				return &fakeAssetService{
					addAssetFunc: func(ctx context.Context, asset models.AddAssetReq) error {
						return nil // simulate success
					},
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "duplicate asset code rejected",
			input: models.AddAssetReq{
				AssetCode:  "AST-001",
				AssetName:  "Dell Laptop",
				CategoryId: 2,
				Status:     "active",
			},
			mockSetup: func() service.AssetService {
				return &fakeAssetService{
					addAssetFunc: func(ctx context.Context, asset models.AddAssetReq) error {
						return errors.New("unique constraint violation")
					},
				}
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := &Handler{Service: tc.mockSetup()} // fresh handler per case, wired to that case's fake service

			body, _ := json.Marshal(tc.input)

			w := httptest.NewRecorder()      // captures what the handler writes, to inspect later
			c, _ := gin.CreateTestContext(w) // fake gin context wrapping the recorder

			c.Request = httptest.NewRequest(http.MethodPost, "/assets", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.AddAsset(c) // call the real handler — it parses the body, calls h.Service.AddAsset, writes the response

			if w.Code != tc.expectedStatus {
				t.Errorf("status = %d; want %d, body: %s", w.Code, tc.expectedStatus, w.Body.String())
			}
		})
	}
}
