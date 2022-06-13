package handlers

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	repo := NewTestRepo()
	NewHandlers(repo)

	os.Exit(m.Run())
}
