package spotify

import (
	"context"
	"fmt"
	"net/http"

	"github.com/andresuryana/y2spot-cli/internal/auth"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(ctx context.Context, auth *auth.Auth) (*Client, error) {
	token, err := auth.LoadToken()
	if err != nil {
		return nil, fmt.Errorf("failed to load token: %w", err)
	}

	tokenSource := auth.Config().TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	if newToken.AccessToken != token.AccessToken {
		_ = auth.SaveToken(newToken)
	}

	return &Client{
		httpClient: auth.Config().Client(ctx, newToken),
	}, nil
}

func (c *Client) GerCurrentUserID(ctx context.Context) (string, error) {
	// TODO: Not yet implemented
	return "user-id", nil
}

func (c *Client) CreatePlaylist(ctx context.Context, userID, name, description string, public bool) (string, error) {
	// TODO: Not yet implemented
	return "", nil
}

func (c *Client) SearchTrackURI(ctx context.Context, artist, title string) (string, error) {
	return "", nil
}

func (c *Client) AddTracksToPlaylist(ctx context.Context, playlistID string, uris []string) error {
	return nil
}
