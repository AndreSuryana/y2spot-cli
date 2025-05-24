package spotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/andresuryana/y2spot-cli/internal/auth"
)

const (
	baseURL = "https://api.spotify.com/v1"

	maxURIs = 100
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

type User struct {
	ID string `json:"id"`
}

func (c *Client) GerCurrentUserID(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/me", nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get current user: %s", resp.Status)
	}

	var data User
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	return data.ID, nil
}

type CreatePlaylistRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type CreatePlaylistResponse struct {
	ID string `json:"id"`
}

func (c *Client) CreatePlaylist(ctx context.Context, userID, name, description string, public bool) (string, error) {
	url := fmt.Sprintf("%s/users/%s/playlists", baseURL, userID)
	body, _ := json.Marshal(&CreatePlaylistRequest{
		Name:        name,
		Description: description,
		Public:      public,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create playlist: %s", resp.Status)
	}

	var playlist CreatePlaylistResponse
	if err := json.NewDecoder(resp.Body).Decode(&playlist); err != nil {
		return "", err
	}

	return playlist.ID, nil
}

type SearchResponse struct {
	Tracks struct {
		Items []struct {
			URI string `json:"uri"`
		} `json:"items"`
	} `json:"tracks"`
}

func (c *Client) SearchTrackURI(ctx context.Context, artist, title string) (string, error) {
	query := url.QueryEscape(fmt.Sprintf("track:%s artist:%s", title, artist))
	url := fmt.Sprintf("%s/search?q=%s&type=track&limit=1", baseURL, query)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to search track: %s", resp.Status)
	}

	var search SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&search); err != nil {
		return "", err
	}

	if len(search.Tracks.Items) == 0 {
		return "", fmt.Errorf("track not found")
	}

	return search.Tracks.Items[0].URI, nil
}

type AddToPlaylistRequest struct {
	URIs []string `json:"uris"`
}

func (c *Client) AddTracksToPlaylist(ctx context.Context, playlistID string, uris []string) (int, error) {
	if len(uris) == 0 {
		return 0, fmt.Errorf("no tracks to add: all searches failed")
	}

	endpoint := fmt.Sprintf("/playlists/%s/tracks", playlistID)
	count := 0

	for i := 0; i < len(uris); i += maxURIs {
		end := min(i+maxURIs, len(uris))

		batch := uris[i:end]
		body, _ := json.Marshal(&AddToPlaylistRequest{URIs: batch})

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+endpoint, bytes.NewReader(body))
		if err != nil {
			return count, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return count, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return count, fmt.Errorf("failed adding track to playlist: %s", resp.Status)
		}

		count += len(batch)
	}

	return count, nil
}
