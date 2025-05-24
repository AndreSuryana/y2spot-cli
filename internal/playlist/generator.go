package playlist

import (
	"context"
	"fmt"

	"github.com/andresuryana/y2spot-cli/internal/domain"
	"github.com/andresuryana/y2spot-cli/internal/spotify"
	"github.com/andresuryana/y2spot-cli/internal/utils"
	"github.com/andresuryana/y2spot-cli/internal/youtube"
)

type Generator struct {
	yt *youtube.Client
	sp *spotify.Client
}

func NewGenerator(ytClient *youtube.Client, spClient *spotify.Client) *Generator {
	return &Generator{
		yt: ytClient,
		sp: spClient,
	}
}

func (g *Generator) GeneratePlaylist(ctx context.Context, playlist *domain.Playlist, ytURL string) (*Result, error) {
	// Fetch the tracks from YouTube Mix/Playlist
	tracks, err := g.yt.FetchTracks(ytURL)
	if err != nil {
		return nil, err
	}

	// Get current user ID
	userID, err := g.sp.GerCurrentUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Spotify user: %w", err)
	}

	// Create Spotify playlist
	playlistID, err := g.sp.CreatePlaylist(
		ctx,
		userID,
		playlist.Name,
		playlist.Description,
		playlist.IsPublic(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create playlist: %w", err)
	}

	var (
		success int
		failed  int
		errors  []string
		uris    []string
	)

	// Add to playlist
	for _, track := range *tracks {
		uri, err := g.sp.SearchTrackURI(ctx, track.Artist, track.Title)
		if err != nil {
			failed++
			errors = append(errors, fmt.Sprintf("%s - %s: %v", track.Artist, track.Title, err))
			continue
		}
		uris = append(uris, uri)
		success++
	}

	// Add to playlist
	err = g.sp.AddTracksToPlaylist(ctx, playlistID, uris)
	if err != nil {
		return nil, fmt.Errorf("failed to add tracks to playlist: %w", err)
	}

	// Write logs
	logPath := utils.WriteErrorLog(errors)

	return &Result{
		NumOfAdded: success,
		NumOfError: failed,
		LogsPath:   logPath,
	}, nil
}
