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

	// these are for testing only
	os.Setenv("SPOTIFY_ID", "d892b49996e1409cb64e9665ae2108e9")
	os.Setenv("SPOTIFY_SECRET", "4255b0a88102430fb49349a5efa38957")

	const redirectURI = "http://localhost:1000/callback"

	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserTopRead, spotifyauth.ScopeUserLibraryModify))
	ch = make(chan *spotify.Client)
	state = "the1"
}

func main() {

	configureServer()

	schedulePlaylistCreation()

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

func schedulePlaylistCreation() {

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

func 