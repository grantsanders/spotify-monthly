package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var (
	auth  *spotifyauth.Authenticator
	ch    chan *spotify.Client
	state string
)

func init() {
// need to set SPOTIFY_ID and SPOTIFY_SECRET variables in docker-compose
	const redirectURI = "http://orb:1000/callback"

	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopePlaylistModifyPublic, spotifyauth.ScopePlaylistModifyPrivate, spotifyauth.ScopeUserTopRead, spotifyauth.ScopeUserLibraryModify))

	ch = make(chan *spotify.Client)
	state = "the1"
}

func main() {

	configureServer()

	client, err := authenticate()
	if err != nil {
		log.Fatal(err)
	}

	schedulePlaylistCreation(client)

	select {}

}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))

	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}

func refresh(client *spotify.Client) *spotify.Client {
	token, err := client.Token()
	if err != nil {
		log.Fatalf("Failed to retrieve existing token")
	}

	newToken, err := auth.RefreshToken(context.Background(), token)
	if err != nil {
		log.Fatalf("Failed to refresh token")
	}

	return spotify.New(auth.Client(context.Background(), newToken))
}

func schedulePlaylistCreation(client *spotify.Client) {

	c := cron.New()
	_, err := c.AddFunc("0 0 1 * *", func() {
		fmt.Printf("Creating playlist at %s", time.Now().Format(time.RFC1123))
		createPlaylist(client)
	})

	if err != nil {
		log.Fatalf("Error scheduling monthly job: %v", err)
	}

	c.Start()
}

func configureServer() {

	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go func() {
		err := http.ListenAndServe(":1000", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

}

func authenticate() (*spotify.Client, error) {

	url := auth.AuthURL(state)

	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	client := <-ch

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("You are logged in as:", user.DisplayName)

	return client, nil

}

func createPlaylist(client *spotify.Client) {

	client = refresh(client)

	now := time.Now().AddDate(0, -1, 0)

	formattedDate := strings.ToLower(now.Format("January '06"))

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	tracks, err := client.CurrentUsersTopTracks(context.Background(), spotify.Timerange(spotify.ShortTermRange), spotify.Limit(30))
	if err != nil {
		log.Fatalf("Error fetching top tracks: %v", err)
	}

	trackIDs := make([]spotify.ID, 0, len(tracks.Tracks))
	for _, track := range tracks.Tracks {
		trackIDs = append(trackIDs, track.ID)
	}

	playlist, err := client.CreatePlaylistForUser(context.Background(), user.ID, formattedDate, "", false, false)
	if err != nil {
		log.Fatalf("Error creating playlist: %v", err)
	}

	_, err = client.AddTracksToPlaylist(context.Background(), playlist.ID, trackIDs...)
	if err != nil {
		log.Fatalf("Error adding tracks to playlist: %v", err)
	}

	log.Println("Playlist created successfully!")
}
