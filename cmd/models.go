package main

import (
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

// Configuration holds our app config
type Config struct {
	PlexURL          string
	AllowedLibraries []string
	RefreshInterval  time.Duration
	CacheDir         string
}

// App represents our application state
type App struct {
	config     Config
	plexServer *PlexServer
	cache      *cache.Cache
	templates  *template.Template
	mu         sync.RWMutex
}

// LibraryItem represents a single item in the Plex library
type LibraryItem struct {
	Title    string    `json:"title"`
	Year     int       `json:"year,omitempty"`
	ThumbURL string    `json:"thumb_url,omitempty"`
	AddedAt  time.Time `json:"added_at"`
	Rating   float64   `json:"rating,omitempty"`
}

// PlexServer handles communication with the Plex server
type PlexServer struct {
	URL    string
	Client *http.Client
}

// PlexResponse represents the XML response from Plex
type PlexResponse struct {
	MediaContainer MediaContainer `xml:"MediaContainer"`
}

type MediaContainer struct {
	Directories []struct {
		Title string `xml:"title,attr"`
		Key   string `xml:"key,attr"`
		Year     int     `xml:"year,attr"`
		Thumb    string  `xml:"thumb,attr,omitempty"`
		AddedAt  int64   `xml:"addedAt,attr"`
		Rating   float64 `xml:"rating,attr,omitempty"`
		ViewGUID string  `xml:"guid,attr"`
	} `xml:"Directory"`
	Videos []struct {
		Title    string  `xml:"title,attr"`
		Year     int     `xml:"year,attr"`
		Thumb    string  `xml:"thumb,attr,omitempty"`
		AddedAt  int64   `xml:"addedAt,attr"`
		Rating   float64 `xml:"rating,attr,omitempty"`
		ViewGUID string  `xml:"guid,attr"`
	} `xml:"Video"`
}
