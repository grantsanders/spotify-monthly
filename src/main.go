package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

// const redirectURI = "localhost:1000/callback"

var (
	auth  *spotifyauth.Authenticator
	ch    chan *spotify.Client
	state string
)

func init() {
	os.Setenv("SPOTIFY_ID", "ID")
	os.Setenv("SPOTIFY_SECRET", "Secret")
	const redirectURI = "http://localhost:1000/callback"

	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserTopRead, spotifyauth.ScopeUserLibraryModify))
	ch = make(chan *spotify.Client)
	state = "abc123"
}

func main() {
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

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	client := <-ch

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	tracks, err := client.CurrentUsersTopTracks(context.Background(), spotify.Timerange(spotify.ShortTermRange), spotify.Limit(30))
	if err != nil {
		log.Fatalf("Error fetching top tracks: %v", err)
	}
	for _, track := range tracks.Tracks {
		fmt.Println(track.Name)
	}

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