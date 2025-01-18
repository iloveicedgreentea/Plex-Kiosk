package main

import (
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
	mu         sync.RWMutex
}

// LibraryItem represents a single item in the Plex library
type LibraryItem struct {
	Title       string    `json:"title"`
	Year        int       `json:"year,omitempty"`
	ThumbURL    string    `json:"thumb_url,omitempty"`
	AddedAt     time.Time `json:"added_at"`
	Rating      float64   `json:"rating,omitempty"`
	Description string    `json:"description,omitempty"`
	Cast        []string  `json:"cast,omitempty"`
	TrailerURL  string    `json:"trailer_url,omitempty"`
}

// PlexServer handles communication with the Plex server
type PlexServer struct {
	URL    string
	Client *http.Client
}

type MediaContainer struct {
	Size      int         `xml:"size,attr"`
	AllowSync int         `xml:"allowSync,attr"`
	Title     string      `xml:"title1,attr"`
	Directory []Directory `xml:"Directory"`
	Videos    []Video     `xml:"Video"`
}

type Directory struct {
	Title     string    `xml:"title,attr"`
	Key       string    `xml:"key,attr"`
	Type      string    `xml:"type,attr"`
	Year      int       `xml:"year,attr"`
	Thumb     string    `xml:"thumb,attr,omitempty"`
	AddedAt   int64     `xml:"addedAt,attr"`
	Rating    float64   `xml:"rating,attr,omitempty"`
	ViewGUID  string    `xml:"guid,attr"`
	Summary   string    `xml:"summary,attr,omitempty"`
	ExtraData ExtraData `xml:"Extras,omitempty"`
	Role      []Role    `xml:"Role,omitempty"`
}

type Video struct {
	Title     string    `xml:"title,attr"`
	Year      int       `xml:"year,attr"`
	Thumb     string    `xml:"thumb,attr,omitempty"`
	AddedAt   int64     `xml:"addedAt,attr"`
	Rating    float64   `xml:"rating,attr,omitempty"`
	ViewGUID  string    `xml:"guid,attr"`
	Summary   string    `xml:"summary,attr,omitempty"`
	ExtraData ExtraData `xml:"Extras,omitempty"`
	Role      []Role    `xml:"Role,omitempty"`
}

type ExtraData struct {
	Size  int     `xml:"size,attr"`
	Extra []Extra `xml:"Video"`
}

type Extra struct {
	Type     string `xml:"type,attr"`
	VideoKey string `xml:"key,attr"`
}

type Role struct {
	Tag string `xml:"tag,attr"`
}
