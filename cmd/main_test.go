package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	config Config
	app    *App
)

func TestMain(m *testing.M) {
	// Run the tests
	os.Setenv("ALLOWED_LIBRARIES", "Movies,TV Shows")

	app = NewApp()
	code := m.Run()
	if code != 0 {
		log.Fatal("Tests failed")
	}
}

func TestFetchLibraryData(t *testing.T) {
	items, err := app.fetchLibraryData()
	require.NoError(t, err)
	require.NotEmpty(t, items["TV Shows"])
	require.NotEmpty(t, items["Movies"])
}
