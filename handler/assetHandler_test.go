package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"inventory-tracker/apperr"
	"inventory-tracker/models"
	service "inventory-tracker/services"
)

func strPtr(s string) *string { return &s } // Go has no address-of for a string constant; wrap it first

// fakeAssetService is a test double for service.AssetService. It lets each case
// decide exactly what the service returns, so the handler can be tested without
// a database, a real service, or a mock of pgx.
//
// Every method is nil-safe: a case that doesn't set a func simply isn't
// exercising that endpoint, and shouldn't panic because of it.
type fakeAssetService struct {
	addAsset       func(ctx context.Context, asset models.AddAssetReq) error
	getAssetList   func(ctx context.Context, params models.AssetListParams) ([]models.AssetResponse, int, error)
	getAssetByCode func(ctx context.Context, assetCode string) (models.AssetResponse, error)
	updateAsset    func(ctx context.Context, id int, asset models.AddAssetReq) error
	deleteByCode   func(ctx context.Context, assetCode string) error
}

func (f *fakeAssetService) AddAsset(ctx context.Context, asset models.AddAssetReq) error {
	if f.addAsset == nil {
		return nil
	}
	return f.addAsset(ctx, asset)
}

func (f *fakeAssetService) GetAssetList(ctx context.Context, params models.AssetListParams) ([]models.AssetResponse, int, error) {
	if f.getAssetList == nil {
		return []models.AssetResponse{}, 0, nil
	}
	return f.getAssetList(ctx, params)
}

func (f *fakeAssetService) GetAssetByCode(ctx context.Context, assetCode string) (models.AssetResponse, error) {
	if f.getAssetByCode == nil {
		return models.AssetResponse{}, nil
	}
	return f.getAssetByCode(ctx, assetCode)
}

func (f *fakeAssetService) UpdateAsset(ctx context.Context, id int, asset models.AddAssetReq) error {
	if f.updateAsset == nil {
		return nil
	}
	return f.updateAsset(ctx, id, asset)
}

func (f *fakeAssetService) DeleteAssetByCode(ctx context.Context, assetCode string) error {
	if f.deleteByCode == nil {
		return nil
	}
	return f.deleteByCode(ctx, assetCode)
}

// compile-time proof the double still satisfies the interface it stands in for.
var _ service.AssetService = (*fakeAssetService)(nil)

func newTestContext(method, target string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, target, bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return c, w
}

func TestAddAsset(t *testing.T) {
	gin.SetMode(gin.TestMode)

	valid := models.AddAssetReq{
		AssetCode:  "AST-001",
		AssetName:  "Dell Laptop",
		Brand:      strPtr("Dell"),
		CategoryId: 2,
		Status:     "active",
		Location:   "HQ",
	}

	cases := []struct {
		name           string
		input          models.AddAssetReq
		svc            service.AssetService
		expectedStatus int
	}{
		{
			// The positive case is what makes the negative ones mean anything: a
			// suite of only rejections passes against code that rejects everything.
			name:           "valid asset is created",
			input:          valid,
			svc:            &fakeAssetService{addAsset: func(context.Context, models.AddAssetReq) error { return nil }},
			expectedStatus: http.StatusCreated,
		},
		{
			name:  "duplicate asset code is a conflict",
			input: valid,
			svc: &fakeAssetService{addAsset: func(context.Context, models.AddAssetReq) error {
				return apperr.ErrDuplicate
			}},
			expectedStatus: http.StatusConflict,
		},
		{
			// Exercises the binding tags: this never reaches the service.
			name:           "missing required fields are rejected before the service",
			input:          models.AddAssetReq{AssetName: "No code"},
			svc:            &fakeAssetService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "an unexpected error is not leaked to the caller",
			input: valid,
			svc: &fakeAssetService{addAsset: func(context.Context, models.AddAssetReq) error {
				return errors.New("connection reset by peer")
			}},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := &Handler{Service: tc.svc}
			body, _ := json.Marshal(tc.input)

			c, w := newTestContext(http.MethodPost, "/assets", body)
			h.AddAsset(c)

			if w.Code != tc.expectedStatus {
				t.Errorf("status = %d; want %d, body: %s", w.Code, tc.expectedStatus, w.Body.String())
			}
			if tc.expectedStatus == http.StatusInternalServerError &&
				bytes.Contains(w.Body.Bytes(), []byte("connection reset")) {
				t.Errorf("internal error text leaked to the client: %s", w.Body.String())
			}
		})
	}
}

func TestDeleteAssetByAssetCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cases := []struct {
		name           string
		svcErr         error
		expectedStatus int
	}{
		{"existing asset is deleted", nil, http.StatusOK},
		{"missing asset is a 404, not a silent 200", apperr.ErrNotFound, http.StatusNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := &Handler{Service: &fakeAssetService{
				deleteByCode: func(context.Context, string) error { return tc.svcErr },
			}}

			c, w := newTestContext(http.MethodDelete, "/assets/AST-001", nil)
			c.Params = gin.Params{{Key: "assetCode", Value: "AST-001"}}
			h.DeleteAssetByAssetCode(c)

			if w.Code != tc.expectedStatus {
				t.Errorf("status = %d; want %d, body: %s", w.Code, tc.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestUpdateAssetRejectsNonNumericID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &Handler{Service: &fakeAssetService{}}
	body, _ := json.Marshal(models.AddAssetReq{
		AssetCode: "AST-001", AssetName: "X", CategoryId: 1, Status: "active", Location: "HQ",
	})

	c, w := newTestContext(http.MethodPut, "/assets/abc", body)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}
	h.UpdateAsset(c)

	// A string where an int belongs is the caller's mistake: 400, not 404 or 500.
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d; want %d, body: %s", w.Code, http.StatusBadRequest, w.Body.String())
	}
}
