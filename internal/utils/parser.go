package utils

import "strings"

func ParseYouTubeTitle(title string) (artist, track string) {
	parts := strings.SplitN(title, " - ", 2)
	if len(parts) != 2 {
		return "", title // fallback
	}

	artist = strings.TrimSpace(parts[0])
	track = strings.TrimSpace(parts[1])

	// Clean suffixes
	cleanSuffixes := []string{
		"(Official Video)", "(Official Video)", "/ Official Video",
	}
	for _, suffix := range cleanSuffixes {
		if strings.HasSuffix(track, suffix) {
			track = strings.TrimSuffix(track, suffix)
			break
		}
	}
	return artist, strings.TrimSpace(track)
}
