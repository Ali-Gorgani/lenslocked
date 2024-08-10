package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Load environment variables from .env file
	dropboxID := os.Getenv("DROPBOX_APP_KEY")
	dropboxSecret := os.Getenv("DROPBOX_APP_SECRET")

	ctx := context.Background()
	config := oauth2.Config{
		ClientID:     dropboxID,
		ClientSecret: dropboxSecret,
		Scopes: []string{
			"files.metadata.read",
			"files.content.read",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.dropbox.com/oauth2/authorize",
			TokenURL: "https://api.dropboxapi.com/oauth2/token",
		},
	}

	url := config.AuthCodeURL(
		"state",
		oauth2.SetAuthURLParam("token_access_type", "offline"),
	)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)
	fmt.Println("Then paste the code here: ")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	client := config.Client(ctx, token)
	resp, err := client.Post(
		"https://api.dropboxapi.com/2/files/list_folder",
		"application/json",
		strings.NewReader(`{"path": ""}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(token)
}
