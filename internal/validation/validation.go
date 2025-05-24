package validation

import (
	"net/url"
	"regexp"
)

func IsYouTubeURL(youTubeURL string) bool {
	parsed, err := url.Parse(youTubeURL)
	if err != nil {
		return false
	}

	// Must be a YouTube domain
	host := parsed.Hostname()
	if host != "www.youtube.com" && host != "youtube.com" && host != "youtu.be" {
		return false
	}

	// Check for playlist parameter or /watch with list param
	query := parsed.Query()
	if listID := query.Get("list"); listID != "" {
		return true
	}

	// For youtu.be or youtube.com/watch, allow only if 'list' param exists
	return false
}

func IsPlaylistNameAllowed(name string) bool {
	// Length check
	if len(name) == 0 || len(name) > 100 {
		return false
	}

	// Basic disallowed characters (you can adjust based on Spotify limits)
	disallowed := regexp.MustCompile(`[<>:"/\\|?*]`)
	return !disallowed.MatchString(name)
}
