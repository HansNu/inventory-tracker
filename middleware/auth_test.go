package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"inventory-tracker/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func signedToken(t *testing.T, method jwt.SigningMethod, key any, claims jwt.Claims) string {
	t.Helper()
	s, err := jwt.NewWithClaims(method, claims).SignedString(key)
	if err != nil {
		t.Fatalf("signing token: %v", err)
	}
	return s
}

func claimsWithExpiry(exp time.Time) utils.Claims {
	return utils.Claims{
		UserID: 42,
		Role:   "IT",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	realSecret := []byte(os.Getenv("JWT_SECRET"))

	validToken, err := utils.GenerateToken(42, "IT")
	if err != nil {
		t.Fatalf("GenerateToken : %v", err)
	}

	noneToken := signedToken(t, jwt.SigningMethodNone,
		jwt.UnsafeAllowNoneSignatureType, claimsWithExpiry(time.Now().Add(time.Hour)))

	wrongKeyToken := signedToken(t, jwt.SigningMethodHS256,
		[]byte("not real secret"), claimsWithExpiry(time.Now().Add(time.Hour)))

	expiredToken := signedToken(t, jwt.SigningMethodHS256,
		realSecret, claimsWithExpiry(time.Now().Add(-time.Hour)))

	cases := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{"no authorization header", "", http.StatusUnauthorized},
		{"missing bearer prefix", validToken, http.StatusUnauthorized},
		{"malformed token", "Bearer not a jwt", http.StatusUnauthorized},
		{"alg none forgery", "Bearer " + noneToken, http.StatusUnauthorized},
		{"signed with wrong key", "Bearer " + wrongKeyToken, http.StatusUnauthorized},
		{"expired Token", "Bearer " + expiredToken, http.StatusUnauthorized},
		{"valid token", "Bearer " + validToken, http.StatusOK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/probe", AuthMiddleware(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"ok": true})
			})

			req := httptest.NewRequest(http.MethodGet, "/probe", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("status = %d; want %d, body: %s",
					w.Code, tc.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestAuthMiddlewareSetsContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	token, err := utils.GenerateToken(42, "IT")
	if err != nil {
		t.Fatalf("GenerateToken : %v", err)
	}

	var gotId, gotRole any

	r := gin.New()
	r.GET("/probe", AuthMiddleware(), func(c *gin.Context) {
		gotId, _ = c.Get("userID")
		gotRole, _ = c.Get("userRole")
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/probe", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if gotId != 42 {
		t.Errorf(" UserId = %v; want 42", gotId)
	}
	if gotRole != "IT" {
		t.Errorf("userRole = %v; want IT", gotRole)
	}
}
