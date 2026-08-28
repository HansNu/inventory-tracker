package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"inventory-tracker/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pashagolub/pgxmock/v4"
)

func strPtr(s string) *string { return &s } //go doesnt have features that allows actual nullable so it would need to serve through pointers

func TestAddAsset(t *testing.T) {
	gin.SetMode(gin.TestMode) //silences gin logging, makes output readable. Doesnt affect behavior

	mockDB, err := pgxmock.NewPool() //fake db connection pool
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	defer mockDB.Close()

	h := &Handler{DB: mockDB}

	// The table. It tells what we send in, what we tell fake DB to do when it receives the query, what http status we expect to return
	cases := []struct {
		name           string
		input          models.AddAssetReq
		mockSetup      func(mockDB pgxmock.PgxPoolIface)
		expectedStatus int
	}{
		{
			name: "valid asset",
			input: models.AddAssetReq{
				AssetCode:    "SCH-01330",
				AssetName:    "Asus Zenbook",
				Brand:        strPtr("Asus"),
				SerialNumber: strPtr("12304990fr"),
				CategoryId:   2, // match whatever type your struct actually declares
				Status:       "Active",
				Location:     "Room 4",
				User:         strPtr("hantest"),
				PurchaseDate: strPtr("2026-07-06"),
				Description:  strPtr("Very good"),
			},
			mockSetup: func(mockDB pgxmock.PgxPoolIface) {
				mockDB.ExpectExec("INSERT INTO asset").
					WithArgs(
						pgxmock.AnyArg(), // asset_code
						pgxmock.AnyArg(), // asset_name
						pgxmock.AnyArg(), // brand
						pgxmock.AnyArg(), // serial_number
						pgxmock.AnyArg(), // category_id
						pgxmock.AnyArg(), // status
						pgxmock.AnyArg(), // location
						pgxmock.AnyArg(), // user
						pgxmock.AnyArg(), // purchase_date
						pgxmock.AnyArg(), // description
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "duplicate asset code rejected",
			input: models.AddAssetReq{
				AssetCode:    "SCH-01330", // same code as above, on purpose
				AssetName:    "Asus Zenbook",
				Brand:        strPtr("Asus"),
				SerialNumber: strPtr("12304990fr"),
				CategoryId:   2,
				Status:       "Active",
				Location:     "Room 4",
				User:         strPtr("hantest"),
				PurchaseDate: strPtr("2026-07-06"),
				Description:  strPtr("Very good"),
			},
			mockSetup: func(mockDB pgxmock.PgxPoolIface) {
				mockDB.ExpectExec("INSERT INTO asset").
					WithArgs(
						pgxmock.AnyArg(), // asset_code
						pgxmock.AnyArg(), // asset_name
						pgxmock.AnyArg(), // brand
						pgxmock.AnyArg(), // serial_number
						pgxmock.AnyArg(), // category_id
						pgxmock.AnyArg(), // status
						pgxmock.AnyArg(), // location
						pgxmock.AnyArg(), // user
						pgxmock.AnyArg(), // purchase_date
						pgxmock.AnyArg(), // description
					).
					WillReturnError(errors.New("unique constraint violation"))
			},
			expectedStatus: http.StatusConflict, // 409 — AddAsset's specific handling for this case
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup(mockDB)

			body, _ := json.Marshal(tc.input)

			w := httptest.NewRecorder() //captures what the handler wrotes to inspect later

			c, _ := gin.CreateTestContext(w) //fake gin context to recorder, same type AddASset Expect as its argument

			c.Request = httptest.NewRequest(http.MethodPost, "/assets", bytes.NewReader(body)) //build real request with json body
			c.Request.Header.Set("Content-Type", "application/json")

			h.AddAsset(c) //call handler

			//check what got written
			if w.Code != tc.expectedStatus {
				t.Errorf("status = %d; want %d, body: %s", w.Code, tc.expectedStatus, w.Body.String())
			}

			//confirm if the mock actually saw the query it expected
			if err := mockDB.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet mock expectations: %v", err)
			}
		})
	}
}
