package handlers

import (
	"os"
	"testing"

	"github.com/Rha02/Carpool-Back-End/config"
)

func TestMain(m *testing.M) {
	var app config.AppConfig

	repo := NewTestRepo(&app)
	NewHandlers(repo)

	os.Exit(m.Run())
}
