package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

// NewApp creates a new application instance
func NewApp(t *template.Template) *App {
	// Load configuration from environment
	config := Config{
		PlexURL:          getEnv("PLEX_URL", "http://localhost:32400"),
		AllowedLibraries: strings.Split(getEnv("ALLOWED_LIBRARIES", ""), ","),
		RefreshInterval:  time.Duration(getIntEnv("REFRESH_INTERVAL", 21600)) * time.Second,
		CacheDir:         "/app/cache",
	}

	if t == nil {
		t = template.Must(template.ParseGlob("templates/*"))
	}
	return &App{
		config: config,
		plexServer: &PlexServer{URL: config.PlexURL, Client: &http.Client{
			Timeout: 10 * time.Second,
		}},
		cache:     cache.New(config.RefreshInterval, config.RefreshInterval*2),
		templates: t,
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getIntEnv gets an environment variable as int with a default value
func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// saveToCacheFile saves the library data to a cache file
func (app *App) saveToCacheFile(data map[string][]LibraryItem) error {
	if err := os.MkdirAll(app.config.CacheDir, 0755); err != nil {
		// assume running locally and try /tmp
		log.Println("Failed to create cache directory, trying /tmp")
		app.config.CacheDir = "/tmp"
		if err := os.MkdirAll(app.config.CacheDir, 0755); err != nil {
			return fmt.Errorf("failed to create cache directory: %w", err)
		}
	}

	cacheFile := filepath.Join(app.config.CacheDir, "library_data.json")
	file, err := os.Create(cacheFile)
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(data)
}

// loadFromCacheFile loads the library data from the cache file
func (app *App) loadFromCacheFile() (map[string][]LibraryItem, error) {
	cacheFile := filepath.Join(app.config.CacheDir, "library_data.json")
	file, err := os.Open(cacheFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open cache file: %w", err)
	}
	defer file.Close()

	var data map[string][]LibraryItem
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode cache file: %w", err)
	}

	return data, nil
}

// refreshCache periodically refreshes the cache
func (app *App) refreshCache() {
	ticker := time.NewTicker(app.config.RefreshInterval)
	for range ticker.C {
		if data, err := app.fetchLibraryData(); err == nil {
			app.mu.Lock()
			app.cache.Set("library_data", data, cache.DefaultExpiration)
			app.mu.Unlock()
			log.Println("Cache refreshed successfully")
		} else {
			log.Printf("Error refreshing cache: %v", err)
		}
	}
}

// setupRoutes configures the HTTP routes
func (app *App) setupRoutes(router *gin.Engine) {
	router.GET("/health", app.healthCheck)
	router.GET("/", app.home)
	router.GET("/thumbnail/*path", app.serveThumb)
}

// healthCheck handles the health check endpoint
func (app *App) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

// serveThumb serves a thumbnail image from Plex to prevent mixed content errors
func (app *App) serveThumb(c *gin.Context) {
	path := c.Param("path")
	// merge plex url with path returned from plex
	fullURL := fmt.Sprintf("%s%s", app.plexServer.URL, path)

	c.Header("Cache-Control", "public, max-age=604800") // 7 days
	c.Header("Expires", time.Now().Add(7*24*time.Hour).UTC().Format(http.TimeFormat))

	// Make the request to Plex
	resp, err := app.plexServer.Client.Get(fullURL)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to fetch thumbnail")
		return
	}
	defer resp.Body.Close()

	// Copy headers from Plex response
	for k, v := range resp.Header {
		c.Header(k, v[0])
	}

	// Stream the response
	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

// home handles the main page
func (app *App) home(c *gin.Context) {
	app.mu.RLock()
	data, found := app.cache.Get("library_data")
	app.mu.RUnlock()

	if !found {
		var err error
		data, err = app.fetchLibraryData()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
		app.mu.Lock()
		app.cache.Set("library_data", data, cache.DefaultExpiration)
		app.mu.Unlock()
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"libraries":   data,
		"lastUpdated": time.Now().Format("2006-01-02 15:04:05"),
	})
}

func (app *App) fetchLibraryData() (map[string][]LibraryItem, error) {
	log.Println("Fetching library data from Plex...")

	libraries, err := app.plexServer.fetchLibraries()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch libraries: %w", err)
	}

	log.Println("Fetched libraries from Plex - ", libraries)

	result := make(map[string][]LibraryItem)
	for _, lib := range libraries {
		// Skip if not in allowed libraries (if any are specified)
		if len(app.config.AllowedLibraries) > 0 {
			found := false
			for _, allowed := range app.config.AllowedLibraries {
				if strings.EqualFold(lib.Title, allowed) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		log.Printf("Processing library: %s", lib.Title)
		items, err := app.plexServer.fetchLibraryContent(lib.Key)
		if err != nil {
			log.Printf("Error fetching content for library %s: %v", lib.Title, err)
			continue
		}
		result[lib.Title] = items
	}

	// Save to cache file
	if err := app.saveToCacheFile(result); err != nil {
		log.Printf("Error saving to cache file: %v", err)
	}

	return result, nil
}

func main() {
	app := NewApp(nil)

	// Start the cache refresh goroutine
	go app.refreshCache()

	// Set up Gin
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.LoadHTMLGlob("templates/*")
	app.setupRoutes(router)

	// Start the server
	log.Fatal(router.Run(":8000"))
}
