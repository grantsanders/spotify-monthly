package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"spotify-monthly/internal/auth"
	"spotify-monthly/internal/playlist"

	"github.com/zmb3/spotify/v2"
)

func ConfigureServer() {
	http.HandleFunc("/callback", completeAuth)
	http.Handle("/ManualCreate", createPlaylistHandler())
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

func completeAuth(w http.ResponseWriter, r *http.Request) {
	state := auth.GetState()

	tok, err := auth.Authenticator.Token(context.Background(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := spotify.New(auth.Authenticator.Client(context.Background(), tok))

	fmt.Fprintf(w, "Login Completed!")
	auth.ClientChannel <- client
}

func createPlaylistHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		playlist.CreatePlaylist()
		fmt.Fprintf(w, "Playlist created successfully")
	}
}
