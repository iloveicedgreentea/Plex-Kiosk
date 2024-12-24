package main

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFetchLib gets all library sections
func TestFetchLibAndData(t *testing.T) {
	ps := PlexServer{
		URL: os.Getenv("PLEX_URL"),
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	lib, err := ps.fetchLibraries()
	require.NoError(t, err)
	assert.NotEmpty(t, lib)
	assert.Len(t, lib, 3)

	for _, l := range lib {
		data, err := ps.fetchLibraryContent(l.Key)
		require.NoError(t, err)
		assert.NotEmpty(t, data)
	}

}
