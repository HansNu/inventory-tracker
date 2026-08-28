package handler

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test-secret-for-tests")
	os.Exit(m.Run())
}
