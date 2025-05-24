package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andresuryana/y2spot-cli/config"
	"github.com/andresuryana/y2spot-cli/internal/auth"
	"github.com/andresuryana/y2spot-cli/internal/domain"
	"github.com/andresuryana/y2spot-cli/internal/playlist"
	"github.com/andresuryana/y2spot-cli/internal/spotify"
	"github.com/andresuryana/y2spot-cli/internal/utils"
	"github.com/andresuryana/y2spot-cli/internal/validation"
	"github.com/andresuryana/y2spot-cli/internal/youtube"
)

func main() {
	// Register interrupt signal handler
	go utils.RegisterSigtermHandler()

	// Load .env files
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("🔐 Config error: %v", err)
	}

	// OAuth2
	id, secret := config.GetSpotifyCredentials()

	if id == "" || secret == "" {
		fmt.Println("🔐 Please provide your Spotify Client ID & Client Secret...")
		fmt.Println("Follow this link: ")
	}

	au := auth.NewAuth(auth.Data{
		ClientID:     id,
		ClientSecret: secret,
		RedirectURL:  "http://127.0.0.1:8080/callback",
		Scopes: []string{
			"playlist-modify-private",
			"playlist-modify-public",
		},
		AuthURL:  "https://accounts.spotify.com/authorize",
		TokenURL: "https://accounts.spotify.com/api/token",
	})

	// Welcome page
	fmt.Println("🎵 Welcome to y2spot-cli - bring your YouTube Mixes & Playlists to Spotify effortlessly!")

	// Authentication checks
	if !au.IsAuthenticated() {
		fmt.Println("\n🔐 Spotify login required. Launching browser for authentication...")
		if err := au.Authenticate(); err != nil {
			fmt.Println("❌ Login failed!", err)
			return
		}

		fmt.Println("✅ Logged in successfully!")
	}

	// YouTube URL (Mix/Playlist)
	var url string
	fmt.Println("\n🔗 Enter the YouTube Mix/Playlist URL:")
	fmt.Print("> ")
	fmt.Scan(&url)

	if !validation.IsYouTubeURL(url) {
		fmt.Println("⚠️ Invalid YouTube Mix/Playlist URL!")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	play := &domain.Playlist{}

	// Spotify playlist name
	fmt.Println("\n🎧 What should we name your new Spotify playlist?")
	fmt.Print("> ")
	play.Name, _ = reader.ReadString('\n')
	play.Name = strings.TrimSpace(play.Name)

	if !validation.IsPlaylistNameAllowed(play.Name) {
		fmt.Println("⚠️ Please use different playlist name, doesn't contains special characters")
		return
	}

	// Spotify playlist name
	fmt.Println("\n🎧 Give you Spotify playlist a brief description?")
	fmt.Print("> ")
	play.Description, _ = reader.ReadString('\n')
	play.Description = strings.TrimSpace(play.Description)

	// Playlist visibility
	var visibility string
	fmt.Println("\n🔒 Should the playlist be public or private? (p = public, v = private)")
	fmt.Print("> ")
	fmt.Scan(&visibility)

	visibility = strings.ToLower(visibility)

	if visibility == "p" {
		play.Visibility = domain.VisibilityPublic
	} else if visibility == "v" {
		play.Visibility = domain.VisibilityPrivate
	} else {
		fmt.Println("⚠️ Invalid playlist visibility!")
		return
	}

	// Confirmation
	fmt.Print("\n🚀 Creating Spotify playlist with:\n\n")
	fmt.Printf("YouTube URL:       %s\n", url)
	fmt.Printf("Playlist Name:     %s\n", play.Name)
	fmt.Printf("Playlist Desc:     %s\n", play.Description)
	fmt.Printf("Visibility:        %s\n", play.Visibility)

	var proc string
	fmt.Println("\nProceeed (Y/n)")
	fmt.Print("> ")
	fmt.Scan(&proc)

	if strings.ToLower(proc) == "y" {
		sp, err := spotify.NewClient(context.Background(), au)
		if err != nil {
			fmt.Println("❌ Failed to create Spotify client", err)
			return
		}

		gn := playlist.NewGenerator(youtube.NewClient(), sp)

		res, err := gn.GeneratePlaylist(context.Background(), play, url)
		if err != nil {
			fmt.Println("❌ Unexpected error", err)
			return
		}

		fmt.Print("\n✅ Created playlist successfully!\n\n")
		fmt.Printf("🎼 %d tracks added\n", res.NumOfAdded)
		fmt.Printf("❌ %d tracks couldn't be found on Spotify\n", res.NumOfError)
		fmt.Printf("\n→ Logs saved in: %s\n", res.LogsPath)
	} else {
		fmt.Println("\n⚠️ Playlist creation canceled")
	}
}
