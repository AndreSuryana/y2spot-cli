package youtube

import (
	"fmt"
	"net/http"

	"github.com/andresuryana/y2spot-cli/internal/domain"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

func (y *Client) FetchTracks(url string) (*domain.Tracks, error) {
	// TODO: Not yet implemented
	return nil, fmt.Errorf("no tracks found")
}
