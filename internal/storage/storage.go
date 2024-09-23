package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
)

func openDBConnection() (*sql.DB, error) {
	connStr := os.Getenv("NEON")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func StoreNewToken(token *oauth2.Token) error {
	db, err := openDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
		UPDATE tokens 
		SET accesstoken = $1, refreshtoken = $2, tokentype = $3, expiry = $4 
		WHERE id = 1`
	_, err = db.Exec(query, token.AccessToken, token.RefreshToken, token.TokenType, token.Expiry)
	if err != nil {
		return fmt.Errorf("failed to update token: %w", err)
	}

	fmt.Println("Token updated successfully")
	return nil
}

func RetrieveLastToken(token *oauth2.Token) (*oauth2.Token, error) {

	db, err := openDBConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT accessToken, tokenType, refreshToken, expiry FROM tokens WHERE id = 1"
	row := db.QueryRow(query)

	err = row.Scan(&token.AccessToken, &token.TokenType, &token.RefreshToken, &token.Expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no token found")
		}
		return nil, fmt.Errorf("failed to retrieve token: %w", err)
	}

	fmt.Println("Retrieved token successfully")

	return token, nil
}
