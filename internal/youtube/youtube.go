package youtube

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/andresuryana/y2spot-cli/config"
	"github.com/andresuryana/y2spot-cli/internal/domain"
	"github.com/andresuryana/y2spot-cli/internal/utils"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Client struct {
	service *youtube.Service
}

func NewClient() (*Client, error) {
	apiKey := config.GetYouTubeAPIKey()

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create YouTube service: %w", err)
	}

	return &Client{service: service}, nil
}

func (y *Client) FetchTracks(ytURL string) (*domain.Tracks, error) {
	parsed, err := url.Parse(ytURL)
	if err != nil {
		return nil, fmt.Errorf("invalid YouTube URL: %w", err)
	}

	query := parsed.Query()
	playlistID := query.Get("list")
	if playlistID == "" {
		return nil, fmt.Errorf("only playlist URLs with 'list' param are supported")
	}

	var tracks domain.Tracks
	pageToken := ""

	// Keep track of the playlist items to prevent duplicate items
	seen := make(map[string]bool)

	for {
		call := y.service.PlaylistItems.List([]string{"snippet"}).
			PlaylistId(playlistID).
			MaxResults(50).
			PageToken(pageToken)

		log.Println("Fetching YouTube playlist...")
		resp, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch playlist items: %w", err)
		}

		for _, item := range resp.Items {
			artist, title := utils.ParseYouTubeTitle(item.Snippet.Title)

			track := &domain.Track{
				Artist: artist,
				Title:  title,
			}

			key := track.Artist + "-" + track.Title
			if !seen[key] {
				tracks = append(tracks, *track)
				seen[key] = true
				log.Printf("Added \"%s - %s\" to track list", track.Artist, track.Title)
			}
		}

		if resp.NextPageToken == "" {
			log.Println("End of YouTube playlist...")
			break
		}
		pageToken = resp.NextPageToken
	}

	if len(tracks) == 0 {
		return nil, fmt.Errorf("no valid tracks found")
	}

	return &tracks, nil
}
