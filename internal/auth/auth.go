package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var (
	Authenticator *spotifyauth.Authenticator
	ClientChannel chan *spotify.Client
	state         string
)

func Setup(redirectUrl string, stateVal string) {
	ClientChannel = make(chan *spotify.Client)
	state = stateVal
	Authenticator = spotifyauth.New(spotifyauth.WithRedirectURL(redirectUrl),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopeUserTopRead,
			spotifyauth.ScopeUserLibraryModify))
}

func GetClient(clientChannel chan *spotify.Client) {

	url := Authenticator.AuthURL(state)

	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	client := <-ClientChannel

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("You are logged in as:", user.DisplayName)

	clientChannel <- client
}

func GetState() string {
	return state
}

func UseRefreshToken(client *spotify.Client) *spotify.Client {
	token, err := client.Token()
	if err != nil {
		log.Println(err.Error())
	}

	newToken, err := Authenticator.RefreshToken(context.Background(), token)
	if err != nil {
		log.Fatalf("Failed to refresh token")
	}

	return spotify.New(Authenticator.Client(context.Background(), newToken))
}
