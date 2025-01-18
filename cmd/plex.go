package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

// fetchLibraries gets all library sections
func (ps *PlexServer) fetchLibraries() ([]struct{ Title, Key string }, error) {
	libraryURL := fmt.Sprintf("%s/library/sections", ps.URL)
	log.Printf("Requesting libraries from: %s", libraryURL)

	resp, err := ps.Client.Get(libraryURL)
	if err != nil {
		log.Printf("Error fetching libraries: %v", err)
		return nil, fmt.Errorf("failed to fetch libraries: %w", err)
	}
	defer resp.Body.Close()

	var plexResp MediaContainer
	decoder := xml.NewDecoder(resp.Body)
	if err := decoder.Decode(&plexResp); err != nil {
		log.Printf("Error decoding XML: %v", err)
		return nil, fmt.Errorf("failed to decode libraries response: %w", err)
	}

	log.Printf("Plex resp %v", plexResp)

	var libraries []struct{ Title, Key string }
	for _, dir := range plexResp.Directory {
		log.Printf("Found library: %s with key: %s", dir.Title, dir.Key)
		libraries = append(libraries, struct{ Title, Key string }{
			Title: dir.Title,
			Key:   dir.Key,
		})
	}
	// TODO: use debug logs
	log.Printf("Total libraries found: %d", len(libraries))
	return libraries, nil
}

func (ps *PlexServer) fetchMetadata(key string) (*MediaContainer, error) {
	resp, err := ps.Client.Get(fmt.Sprintf("%s/library/metadata/%s", ps.URL, key))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata body: %w", err)
	}

	var plexResp MediaContainer
	if err := xml.Unmarshal(body, &plexResp); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	return &plexResp, nil
}

// getTrailerURL attempts to find a trailer URL from extras
func (ps *PlexServer) getTrailerURL(extras ExtraData) string {
	for _, extra := range extras.Extra {
		if strings.EqualFold(extra.Type, "trailer") {
			return fmt.Sprintf("/thumbnail%s", extra.VideoKey)
		}
	}
	return ""
}

// fetchLibraryContent gets all items in a library section
func (ps *PlexServer) fetchLibraryContent(key string) ([]LibraryItem, error) {
	log.Println("using key", key)
	resp, err := ps.Client.Get(fmt.Sprintf("%s/library/sections/%s/all", ps.URL, key))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch library content: %w", err)
	}
	defer resp.Body.Close()

	var plexResp MediaContainer
	if err := xml.NewDecoder(resp.Body).Decode(&plexResp); err != nil {
		return nil, fmt.Errorf("failed to decode library content: %w", err)
	}

	var items []LibraryItem

	// Handle video content (movies)
	for _, video := range plexResp.Videos {
		thumb := ""
		if video.Thumb != "" {
			thumb = video.Thumb
		}

		var cast []string
		for _, role := range video.Role {
			cast = append(cast, role.Tag)
		}

		items = append(items, LibraryItem{
			Title:       video.Title,
			Year:        video.Year,
			ThumbURL:    fmt.Sprintf("/thumbnail%s", thumb),
			AddedAt:     time.Unix(video.AddedAt, 0),
			Rating:      video.Rating,
			Description: video.Summary,
			Cast:        cast,
			TrailerURL:  ps.getTrailerURL(video.ExtraData),
		})
	}

	// Handle directory content (TV shows)
	for _, dir := range plexResp.Directory {
		thumb := ""
		if dir.Thumb != "" {
			thumb = dir.Thumb
		}

		var cast []string
		for _, role := range dir.Role {
			cast = append(cast, role.Tag)
		}

		items = append(items, LibraryItem{
			Title:       dir.Title,
			Year:        dir.Year,
			ThumbURL:    fmt.Sprintf("/thumbnail%s", thumb),
			AddedAt:     time.Unix(dir.AddedAt, 0),
			Rating:      dir.Rating,
			Description: dir.Summary,
			Cast:        cast,
			TrailerURL:  ps.getTrailerURL(dir.ExtraData),
		})
	}

	return items, nil
}
