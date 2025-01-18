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
	resp, err := ps.Client.Get(fmt.Sprintf("%s/library/sections", ps.URL))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch libraries: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var plexResp MediaContainer
	if err := xml.Unmarshal(body, &plexResp); err != nil {
		return nil, fmt.Errorf("failed to decode libraries response: %w", err)
	}

	var libraries []struct{ Title, Key string }

	for _, dir := range plexResp.Directories {
		libraries = append(libraries, struct{ Title, Key string }{
			Title: dir.Title,
			Key:   strings.TrimSpace(dir.Key),
		})
	}
	return libraries, nil
}

// fetchLibraryContent gets all items in a library section
func (ps *PlexServer) fetchLibraryContent(key string) ([]LibraryItem, error) {
	log.Println("using key", key)
	resp, err := ps.Client.Get(fmt.Sprintf("%s/library/sections/%s/all", ps.URL, key))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch library content: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var plexResp MediaContainer
	if err := xml.Unmarshal(body, &plexResp); err != nil {
		return nil, fmt.Errorf("failed to decode library content: %w", err)
	}

	var items []LibraryItem

	// plex stupidly uses two different structs and fields for tv/movie
	// so we have to check which one we're dealing with
	isVideo := len(plexResp.Videos) > 0
	isDirectory := len(plexResp.Directories) > 0

	// TODO: use generic functions to DRY
	if isVideo {
		for _, video := range plexResp.Videos {
			thumb := ""
			if video.Thumb != "" {
				thumb = video.Thumb
			}

			items = append(items, LibraryItem{
				Title:    video.Title,
				Year:     video.Year,
				ThumbURL: fmt.Sprintf("/thumbnail%s", thumb), // using nginx to serve the content
				AddedAt:  time.Unix(video.AddedAt, 0),
				Rating:   video.Rating,
			})
		}
	} else if isDirectory {
		for _, video := range plexResp.Directories {
			thumb := ""
			if video.Thumb != "" {
				thumb = video.Thumb
			}

			items = append(items, LibraryItem{
				Title:    video.Title,
				Year:     video.Year,
				ThumbURL: fmt.Sprintf("/thumbnail%s", thumb),
				AddedAt:  time.Unix(video.AddedAt, 0),
				Rating:   video.Rating,
			})
		}
	}

	return items, nil
}
