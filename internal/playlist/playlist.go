package playlist

import (
	"context"
	"log"
	"spotify-monthly/internal/auth"
	"strings"
	"time"

	"github.com/zmb3/spotify/v2"
)

var (
	ClientChannel chan *spotify.Client = make(chan *spotify.Client, 1)
)

func CreatePlaylist() {

	client := <-ClientChannel

	client = useRefreshToken(client)

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

	ClientChannel <- client
	log.Println("Playlist created successfully!")
}

func useRefreshToken(client *spotify.Client) *spotify.Client {
	token, err := client.Token()
	if err != nil {
		log.Println(err.Error())
	}

	newToken, err := auth.Authenticator.RefreshToken(context.Background(), token)
	if err != nil {
		log.Fatalf("Failed to refresh token")
	}

	return spotify.New(auth.Authenticator.Client(context.Background(), newToken))
}
