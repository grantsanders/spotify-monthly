package auth

import (
	"context"
	"fmt"
	"log"
	"spotify-monthly/internal/storage"
	"time"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
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

	fmt.Println(redirectUrl)
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

func GetTokenFromDB() (*spotify.Client, error) {

	token := &oauth2.Token{
		AccessToken:  "",
		TokenType:    "",
		RefreshToken: "",
		Expiry:       time.Now().Add(1 * time.Hour),
	}

	token, err := storage.RetrieveLastToken(token)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	newToken, err := Authenticator.RefreshToken(context.Background(), token)
	if err != nil {
		fmt.Println("Failed to refresh token")
		return nil, err
	}

	return spotify.New(Authenticator.Client(context.Background(), newToken)), nil
}
